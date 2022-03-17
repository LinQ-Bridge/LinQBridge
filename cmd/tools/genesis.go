package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/urfave/cli"
	"gorm.io/gorm"

	"land-bridge/models"
	networkUtils "land-bridge/network/utils"
	"land-bridge/utils"
)

var pubkeyList string

var GenesisCMD = cli.Command{
	Name:  "genesis",
	Usage: "You must point to a JSON file that holds all chain public key addresses",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "pklist",
			Usage:       "Node Pubkey Address set JSON",
			Required:    true,
			Destination: &pubkeyList,
		},
	},
	Action: genesis,
}

func genesis(ctx *cli.Context) {
	pkstr := strings.Split(pubkeyList, ",")

	var addrs []common.Address
	for _, s := range pkstr {
		pub := common.HexToAddress(s)
		addrs = append(addrs, pub)
	}

	extraData, err := Encode("0x00", addrs)
	if err != nil {
		logs.Error("Failed to encode extra data", "err", err)
		return
	}

	time := time.Now().Unix()
	block := new(networkUtils.Block)
	block.ParentHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	block.Coinbase = common.HexToAddress("0x0000000000000000000000000000000000000000")
	block.Time = uint64(time)
	block.Height = 0
	block.Type = 1
	block.SrcTx = networkUtils.TxInfo{
		ChainID: 0,
		TxHash:  []byte("0x0000000000000000000000000000000000000000000000000000000000000000"),
	}
	block.DstTx = networkUtils.TxInfo{
		ChainID: 0,
		TxHash:  []byte("0x0000000000000000000000000000000000000000000000000000000000000000"),
	}
	block.ExtraData = common.Hex2Bytes(extraData[2:])
	block.Hash()

	fileStr := fmt.Sprintf(GenesisTemplate, block.BlockHash, strconv.FormatInt(time, 16), extraData)

	fileName := "genesis.json"
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		return
	}

	_, err = f.Write([]byte(fileStr))
	if err != nil {
		return
	}
}

func Encode(vanity string, validators []common.Address) (string, error) {
	newVanity, err := hexutil.Decode(vanity)
	if err != nil {
		return "", err
	}

	if len(newVanity) < ExtraVanity {
		newVanity = append(newVanity, bytes.Repeat([]byte{0x00}, ExtraVanity-len(newVanity))...)
	}
	newVanity = newVanity[:ExtraVanity]

	ist := &networkUtils.LBFTExtra{
		Validators:    validators,
		Seal:          make([]byte, 0),
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + common.Bytes2Hex(append(newVanity, payload...)), nil
}

const (
	ExtraVanity     = 32
	GenesisTemplate = `{
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "blockHash": "%s",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "height": "0x0",
  "timestamp": "0x%s",
  "type": "0x1",
  "srcTx": {
    "chain_id": "0x0",
    "tx_hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
  },
  "cBytes": "0x00000000000000000000000000000000000000000000",
  "dstTx": {
    "chain_id": "0x0",
    "tx_hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
  },
  "extraData": "%s"
}`
)

func initGenesis(genesisPath string, db *gorm.DB) error {
	file, err := os.Open(genesisPath)
	if err != nil {
		utils.Fatalf("Failed to read genesis file: %v", err)
	}
	defer file.Close()

	var (
		genesis = networkUtils.Block{}
		model   = &models.Block{}
		buf     = new(bytes.Buffer)
	)

	if err := json.NewDecoder(file).Decode(&genesis); err != nil {
		utils.Fatalf("invalid genesis file: %v", err)
	}

	err = rlp.Encode(buf, &genesis)
	if err != nil {
		return err
	}
	model.Bytes = common.Bytes2Hex(buf.Bytes())
	model.Height = genesis.Height
	model.BlockHash = genesis.BlockHash.Hex()
	if err := db.Create(model).Error; err != nil {
		return err
	}

	return nil
}
