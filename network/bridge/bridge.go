package bridge

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"land-bridge/conf"
	"land-bridge/constant"
	"land-bridge/contracts/eccm"
	"land-bridge/models"
)

type Bridge struct {
	db         *gorm.DB
	bq         *Queryer
	proxyAddrs map[uint64]string
	signer     *Signer
	priv       *ecdsa.PrivateKey
	chainMap   map[uint64]*conf.ChainListenConfig
}

func NewBridge(db *gorm.DB, cfg *conf.Config, priv *ecdsa.PrivateKey) *Bridge {
	mapProxyAddrs := make(map[uint64]string)

	for _, chain := range cfg.Chains {
		mapProxyAddrs[chain.ChainID] = chain.NFTProxyContract
	}

	addr := crypto.PubkeyToAddress(priv.PublicKey)
	signer := NewSigner(priv, &addr)
	bridgeQueryer := NewBridgeQueryer(cfg)

	chainMap := make(map[uint64]*conf.ChainListenConfig)
	for _, chain := range cfg.Chains {
		chainMap[chain.ChainID] = chain
	}

	return &Bridge{
		db:         db,
		bq:         bridgeQueryer,
		proxyAddrs: mapProxyAddrs,
		signer:     signer,
		priv:       priv,
		chainMap:   chainMap,
	}
}

func (b *Bridge) Sign(tx *TxParam) ([]byte, error) {
	return b.signer.Sign(tx)
}

func (b *Bridge) UpdateWrapper(txhash common.Hash) error {
	hashStr := txhash.Hex()[2:]
	wrapperTransaction := new(models.WrapperTransaction)
	if err := b.db.Model(&wrapperTransaction).Where("hash = ?", hashStr).Update("status", constant.STATE_SOURCE_CONFIRMED).Error; err != nil {
		logs.Error("GORM UpdateWrapper err", err)
		return err
	}

	return nil
}

func (b *Bridge) PendingWrapper(wt *models.WrapperTransaction) error {
	b.db.Model(wt).Update("status", constant.STATE_PENDDING)
	return nil
}

func (b *Bridge) PendingWrapperSkip(wt *models.WrapperTransaction) error {
	b.db.Model(wt).Update("status", constant.STATE_SOURCE_CONFIRMED)
	return nil
}

func (b *Bridge) BridgeMakeTx(wrapperTransaction *models.WrapperTransaction) (*TxParam, error) {
	srcTransfer := new(models.SrcTransfer)
	if err := b.db.Where("tx_hash = ?", wrapperTransaction.Hash).First(&srcTransfer).Error; err != nil {
		return nil, err
	}
	tokenURI, err := b.bq.GetTokenURIByAssetWithID(wrapperTransaction.SrcChainID, srcTransfer.Asset, &srcTransfer.TokenID.Int)
	if err != nil {
		logs.Error("GetTokenURIByAssetWithID error", err)
		return nil, err
	}

	tx := ConstructTx(wrapperTransaction, srcTransfer, b.proxyAddrs[wrapperTransaction.DstChainID], tokenURI)

	return tx, nil
}

func (b *Bridge) BridgeToChainB(txhash common.Hash, signatures map[common.Address][]byte) error {
	hashStr := txhash.Hex()[2:]
	wrapperTransaction := new(models.WrapperTransaction)
	if err := b.db.Where("hash = ?", hashStr).First(&wrapperTransaction).Error; err != nil {
		logs.Error("BridgeToChainB wrapperTransaction", err)
		return err
	}
	srcTransfer := new(models.SrcTransfer)
	if err := b.db.Where("tx_hash = ?", hashStr).First(&srcTransfer).Error; err != nil {
		logs.Error("BridgeToChainB srcTransfer", err)
		return err
	}

	tokenURI, err := b.bq.GetTokenURIByAssetWithID(wrapperTransaction.SrcChainID, srcTransfer.Asset, &srcTransfer.TokenID.Int)
	if err != nil {
		logs.Error("GetTokenURIByAssetWithID error", err)
		return err
	}

	tx := ConstructTx(wrapperTransaction, srcTransfer, b.proxyAddrs[wrapperTransaction.DstChainID], tokenURI)

	b.db.Model(wrapperTransaction).Update("status", constant.STATE_SOURCE_CONFIRMED)

	var argSignature []byte
	for _, s := range signatures {
		argSignature = append(argSignature, s...)
	}
	tx.SetSignatures(argSignature)

	chainConf := b.chainMap[tx.ChainID()]
	var rawClient *ethclient.Client
	for rawClient == nil {
		for _, s := range chainConf.GetNodesURL() {
			rawClient, _ = ethclient.Dial(s)
			if rawClient != nil {
				break
			}
		}
	}
	execTxHash, err := b.transactionExec(tx, chainConf, rawClient)
	if err != nil {
		logs.Error("transactionExec error", err)
		errorT := &models.ErrorTransaction{
			TxHash:       wrapperTransaction.Hash,
			FromChainID:  wrapperTransaction.SrcChainID,
			FromContract: srcTransfer.Asset,
			ToChainID:    wrapperTransaction.DstChainID,
			ToContract:   b.proxyAddrs[wrapperTransaction.DstChainID],
			ToAssetHash:  srcTransfer.DstAsset,
			ToAddress:    srcTransfer.DstUser,
			TokenID:      srcTransfer.TokenID,
			TokenURI:     tokenURI,
			Signature:    common.Bytes2Hex(argSignature),
			ErrorMsg:     err.Error(),
		}
		b.db.Create(errorT)
		return err
	}

	logs.Info("bridge cross txHash:", execTxHash.Hex())

	return nil
}

func (b *Bridge) transactionExec(tx *TxParam, chainConf *conf.ChainListenConfig, client *ethclient.Client) (common.Hash, error) {
	ccmContractAddr := common.HexToAddress(chainConf.CCMContract)
	ccm, err := eccm.NewEthCrossChainManager(ccmContractAddr, client)
	if err != nil {
		return common.Hash{}, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(b.priv, big.NewInt(int64(chainConf.ChainID)))
	if err != nil {
		return common.Hash{}, err
	}
	auth.GasLimit = uint64(400000)

	// TODO Test code - the lowest gasPrice
	gasPrice := big.NewInt(0)
	switch chainConf.ChainID {
	case 97:
		gasPrice = big.NewInt(10000000000)
	case 1001:
		gasPrice = big.NewInt(750000000000)
	case 210309:
		gasPrice = big.NewInt(1000000000)
	default:
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			return common.Hash{}, err
		}
	}
	auth.GasPrice = gasPrice
	// ------------

	txback, err := ccm.VerifySigAndExecuteTx(auth, tx.Serialize(), tx.GetSignatures())
	if err != nil {
		return common.Hash{}, err
	}

	return txback.Hash(), nil
}
