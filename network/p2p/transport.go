package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/bitutil"
	"github.com/ethereum/go-ethereum/rlp"

	"land-bridge/network/p2p/rlpx"
)

const (
	// total timeout for encryption handshake and protocol
	// handshake in both directions.
	handshakeTimeout = 5 * time.Second

	// This is the timeout for sending the disconnect reason.
	// This is shorter than the usual timeout because we don't want
	// to wait if the connection is known to be bad anyway.
	discWriteTimeout = 1 * time.Second
)

type transport interface {
	// The two handshakes.
	doEncHandshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error)

	doProtoHandshake(our *protoHandshake) (*protoHandshake, error)

	// The MsgReadWriter can only be used after the encryption
	// handshake has completed. The code uses conn.id to track this
	// by setting it to a non-nil value after the encryption handshake.
	MsgReadWriter
	// transports must provide Close because we use MsgPipe in some of
	// the tests. Closing the actual network connection doesn't do
	// anything in those tests because MsgPipe doesn't use it.
	close(err error)
}

// rlpxTransport is the transport used by actual (non-test) connections.
// It wraps an RLPx connection with locks and read/write deadlines.
type rlpxTransport struct {
	rmu, wmu sync.Mutex
	wbuf     bytes.Buffer
	conn     *rlpx.Conn
}

func newRLPX(conn net.Conn, dialDest *ecdsa.PublicKey) transport {
	return &rlpxTransport{conn: rlpx.NewConn(conn, dialDest)}
}

func (t *rlpxTransport) doEncHandshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	t.conn.SetDeadline(time.Now().Add(handshakeTimeout))
	return t.conn.Handshake(prv)
}

func (t *rlpxTransport) doProtoHandshake(our *protoHandshake) (their *protoHandshake, err error) {
	// Writing our handshake happens concurrently, we prefer
	// returning the handshake read error. If the remote side
	// disconnects us early with a valid reason, we should return it
	// as the error so it can be tracked elsewhere.
	werr := make(chan error, 1)
	go func() { werr <- Send(t, handshakeMsg, our) }()
	if their, err = readProtocolHandshake(t); err != nil {
		<-werr // make sure the write terminates too
		return nil, err
	}
	if err := <-werr; err != nil {
		return nil, fmt.Errorf("write error: %v", err)
	}
	// If the protocol version supports Snappy encoding, upgrade immediately
	t.conn.SetSnappy(their.Version >= snappyProtocolVersion)

	return their, nil
}

func (t *rlpxTransport) ReadMsg() (Msg, error) {
	t.rmu.Lock()
	defer t.rmu.Unlock()

	var msg Msg
	t.conn.SetReadDeadline(time.Now().Add(frameReadTimeout))
	code, data, wireSize, err := t.conn.Read()
	if err == nil {
		// Protocol messages are dispatched to subprotocol handlers asynchronously,
		// but package rlpx may reuse the returned 'data' buffer on the next call
		// to Read. Copy the message data to avoid this being an issue.
		data = common.CopyBytes(data)
		msg = Msg{
			ReceivedAt: time.Now(),
			Code:       code,
			Size:       uint32(len(data)),
			meterSize:  uint32(wireSize),
			Payload:    bytes.NewReader(data),
		}
	}
	return msg, err
}

func (t *rlpxTransport) WriteMsg(msg Msg) error {
	t.wmu.Lock()
	defer t.wmu.Unlock()

	// Copy message data to write buffer.
	t.wbuf.Reset()
	if _, err := io.CopyN(&t.wbuf, msg.Payload, int64(msg.Size)); err != nil {
		return err
	}

	// Write the message.
	t.conn.SetWriteDeadline(time.Now().Add(frameWriteTimeout))
	size, err := t.conn.Write(msg.Code, t.wbuf.Bytes())
	if err != nil {
		return err
	}
	msg.meterSize = size

	return nil
}

func (t *rlpxTransport) close(err error) {
	t.wmu.Lock()
	defer t.wmu.Unlock()

	// Tell the remote end why we're disconnecting if possible.
	// We only bother doing this if the underlying connection supports
	// setting a timeout tough.
	if t.conn != nil {
		if r, ok := err.(DiscReason); ok && r != DiscNetworkError {
			deadline := time.Now().Add(discWriteTimeout)
			if err := t.conn.SetWriteDeadline(deadline); err == nil {
				// Connection supports write deadline.
				t.wbuf.Reset()
				rlp.Encode(&t.wbuf, []DiscReason{r})
				t.conn.Write(discMsg, t.wbuf.Bytes())
			}
		}
	}
	t.conn.Close()
}

func readProtocolHandshake(rw MsgReader) (*protoHandshake, error) {
	msg, err := rw.ReadMsg()
	if err != nil {
		return nil, err
	}
	if msg.Size > baseProtocolMaxMsgSize {
		return nil, fmt.Errorf("message too big")
	}
	if msg.Code == discMsg {
		// Disconnect before protocol handshake is valid according to the
		// spec and we send it ourself if the post-handshake checks fail.
		// We can't return the reason directly, though, because it is echoed
		// back otherwise. Wrap it in a string instead.
		var reason [1]DiscReason
		rlp.Decode(msg.Payload, &reason)
		return nil, reason[0]
	}
	if msg.Code != handshakeMsg {
		return nil, fmt.Errorf("expected handshake, got %x", msg.Code)
	}
	var hs protoHandshake
	if err := msg.Decode(&hs); err != nil {
		return nil, err
	}
	if len(hs.ID) != 64 || !bitutil.TestBytes(hs.ID) {
		return nil, DiscInvalidIdentity
	}
	return &hs, nil
}
