package node

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/event"
	"gorm.io/gorm"

	"land-bridge/network/bridge"
	"land-bridge/network/linq/txblock"
	"land-bridge/network/p2p"
)

const (
	initializingState = iota
	runningState
	closedState
)

type Node struct {
	eventmux      *event.TypeMux // Event multiplexer used between the services of a stack
	config        *Config
	server        *p2p.Server
	stop          chan struct{}
	startStopLock sync.Mutex  // Start/Stop are protected by an additional lock
	state         int         // Tracks state of node lifecycle
	lifecycles    []Lifecycle // All registered backends, services, and auxiliary services that have a lifecycle
	lock          sync.Mutex

	db     *gorm.DB
	Bridge *bridge.Bridge
	Pool   *txblock.BlockPool
}

func NewNode(conf *Config) (*Node, error) {
	confCopy := *conf
	conf = &confCopy
	if conf.DataDir != "" {
		absdatadir, err := filepath.Abs(conf.DataDir)
		if err != nil {
			return nil, err
		}
		conf.DataDir = absdatadir
	}

	node := &Node{
		config:   conf,
		stop:     make(chan struct{}),
		server:   &p2p.Server{Config: conf.P2P},
		eventmux: new(event.TypeMux),
	}
	node.server.Config.PrivateKey = node.config.NodeKey()
	node.server.Config.Name = node.config.NodeName()

	if node.server.Config.StaticNodes == nil {
		node.server.Config.StaticNodes = conf.P2P.BootstrapNodes
	}
	return node, nil
}

func (n *Node) Start() error {
	n.startStopLock.Lock()
	defer n.startStopLock.Unlock()

	n.lock.Lock()
	switch n.state {
	case runningState:
		n.lock.Unlock()
		return errors.New("node already running")
	case closedState:
		n.lock.Unlock()
		return errors.New("node not started")
	}
	n.state = runningState

	logs.Info("Starting peer-to-peer node", "instance", n.server.Name)
	err := n.server.Start()
	if err != nil {
		logs.Error("peer-to-peer node err:", err)
		n.server.Stop()
		n.doClose(nil)
	}

	lifecycles := make([]Lifecycle, len(n.lifecycles))
	copy(lifecycles, n.lifecycles)
	n.lock.Unlock()

	logs.Info("self node enode:", n.server.NodeInfo().Enode)

	// Start all registered lifecycles.
	var started []Lifecycle
	for _, lifecycle := range lifecycles {
		if err = lifecycle.Start(); err != nil {
			break
		}
		started = append(started, lifecycle)
	}

	logs.Info("server start lifecycles ok")

	if err != nil {
		n.stopServices(started)
		n.doClose(nil)
	}

	return err
}

func (n *Node) stopServices(running []Lifecycle) error {
	// Stop running lifecycles in reverse order.
	failure := &StopError{Services: make(map[reflect.Type]error)}
	for i := len(running) - 1; i >= 0; i-- {
		if err := running[i].Stop(); err != nil {
			failure.Services[reflect.TypeOf(running[i])] = err
		}
	}

	// Stop p2p networking.
	n.server.Stop()

	if len(failure.Services) > 0 {
		return failure
	}
	return nil
}

func (n *Node) SetDB(db *gorm.DB) {
	n.db = db
}

func (n *Node) GetDB() *gorm.DB {
	return n.db
}

func (n *Node) Close() error {
	n.startStopLock.Lock()
	defer n.startStopLock.Unlock()

	n.lock.Lock()
	state := n.state
	n.lock.Unlock()

	switch state {
	case initializingState, runningState:
		// The node was never started.
		return n.doClose(nil)
	case closedState:
		return errors.New("node not started")
	default:
		panic(fmt.Sprintf("node is in unknown state %d", state))
	}
}

// doClose releases resources acquired by New(), collecting errors.
func (n *Node) doClose(errs []error) error {
	// Close databases. This needs the lock because it needs to
	// synchronize with OpenDatabase*.
	n.lock.Lock()
	n.state = closedState
	n.lock.Unlock()

	// Unblock n.Wait.
	close(n.stop)

	// Report any errors that might have occurred.
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return fmt.Errorf("%v", errs)
	}
}

// Wait blocks until the node is closed.
func (n *Node) Wait() {
	<-n.stop
}

// Server retrieves the currently running P2P network layer. This method is meant
// only to inspect fields of the currently running server. Callers should not
// start or stop the returned server.
func (n *Node) Server() *p2p.Server {
	n.lock.Lock()
	defer n.lock.Unlock()

	return n.server
}

// RegisterProtocols adds backend's protocols to the node's p2p server.
func (n *Node) RegisterProtocols(protocols []p2p.Protocol) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.state != initializingState {
		panic("can't register protocols on running/stopped node")
	}
	n.server.Protocols = append(n.server.Protocols, protocols...)
}

// RegisterLifecycle registers the given Lifecycle on the node.
func (n *Node) RegisterLifecycle(lifecycle Lifecycle) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.state != initializingState {
		panic("can't register lifecycle on running/stopped node")
	}
	if containsLifecycle(n.lifecycles, lifecycle) {
		panic(fmt.Sprintf("attempt to register lifecycle %T more than once", lifecycle))
	}
	n.lifecycles = append(n.lifecycles, lifecycle)
}

// containsLifecycle checks if 'lfs' contains 'l'.
func containsLifecycle(lfs []Lifecycle, l Lifecycle) bool {
	for _, obj := range lfs {
		if obj == l {
			return true
		}
	}
	return false
}

func (n *Node) GetNodeKey() *ecdsa.PrivateKey {
	return n.config.NodeKey()
}

// EventMux retrieves the event multiplexer used by all the network services in
// the current protocol stack.
func (n *Node) EventMux() *event.TypeMux {
	return n.eventmux
}
