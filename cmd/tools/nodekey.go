package tools

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

var NodekeyCMD = cli.Command{
	Name:   "nodekey",
	Usage:  "Get a random nodekey file and its enode information",
	Action: nodekey,
}

func nodekey(ctx *cli.Context) {
	// No persistent key found, generate and store a new one.
	key, err := crypto.GenerateKey()
	if err != nil {
		logs.Warn(fmt.Sprintf("Failed to generate node key: %v", err))
	}

	pwd, _ := os.Getwd()

	keyfile := filepath.Join(pwd, "nodekey")
	if err := crypto.SaveECDSA(keyfile, key); err != nil {
		logs.Error(fmt.Sprintf("Failed to persist node key: %v", err))
	}

	pubkey := key.Public().(*ecdsa.PublicKey)
	nodeid := fmt.Sprintf("%x", crypto.FromECDSAPub(pubkey)[1:])

	enode := fmt.Sprintf(enodeTemplate, nodeid)

	publicAddr := crypto.PubkeyToAddress(*pubkey)

	fmt.Println("linq NodeKey successfully generated")
	fmt.Println("The new file is generated in the current directory")
	fmt.Println("Pubkey Address:" + publicAddr.Hex())
	fmt.Println(enode)
}

const (
	enodeTemplate = "enode://%s@127.0.0.1:30303"
)
