package platonlisten

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"land-bridge/conf"
	"land-bridge/constant"
	"land-bridge/contracts/eccm"
	"land-bridge/contracts/nftlp"
	"land-bridge/contracts/nftwrap"
	"land-bridge/handle/chainclient"
	"land-bridge/handle/listener/utils"
	"land-bridge/models"
)

type PlatonChainListen struct {
	platonCfg *conf.ChainListenConfig
	platonSdk *chainclient.PlatonSdkPro
}

func NewPlatonChainListen(cfg *conf.ChainListenConfig) *PlatonChainListen {
	urls := cfg.GetNodesURL()
	sdk := chainclient.NewPlatonSdkPro(urls, cfg.ListenSlot, cfg.ChainID)
	listen := &PlatonChainListen{cfg, sdk}
	return listen
}

func (g *PlatonChainListen) GetChainName() string {
	return g.platonCfg.ChainName
}

func (g *PlatonChainListen) GetChainID() uint64 {
	return g.platonCfg.ChainID
}

func (g *PlatonChainListen) GetChainListenSlot() uint64 {
	return g.platonCfg.ListenSlot
}

func (g *PlatonChainListen) GetBatchSize() uint64 {
	return g.platonCfg.BatchSize
}

func (g *PlatonChainListen) GetDefer() uint64 {
	return g.platonCfg.Defer
}

func (g *PlatonChainListen) GetLatestHeight() (uint64, error) {
	return g.platonSdk.GetLatestHeight()
}

func (g *PlatonChainListen) HandleNewBlock(height uint64) ([]*models.WrapperTransaction, []*models.SrcTransaction, []*models.DstTransaction, int, int, error) {
	header, err := g.platonSdk.GetHeaderByNumber(height)
	if err != nil {
		return nil, nil, nil, 0, 0, err
	}
	if header == nil {
		return nil, nil, nil, 0, 0, fmt.Errorf("there is no geth block")
	}
	time := header.Time
	wrapperTransactions := make([]*models.WrapperTransaction, 0)
	nftWrapperTransactions, err := g.getNFTWrapperEventByBlockNumber(g.platonCfg.NFTWrapperContract, height, height)
	if err != nil {
		return nil, nil, nil, 0, 0, err
	}
	wrapperTransactions = append(wrapperTransactions, nftWrapperTransactions...)

	for _, item := range wrapperTransactions {
		//logs.Info("(wrapper) from chain: %s, txhash: %s", g.GetChainName(), item.Hash)
		item.Time = time
		item.SrcChainID = g.GetChainID()
		item.Status = constant.STATE_SOURCE_DONE
	}

	eccmLockEvents, eccmUnLockEvents, err := g.getECCMEventByBlockNumber(g.platonCfg.CCMContract, height, height)
	if err != nil {
		return nil, nil, nil, 0, 0, err
	}

	proxyLockEvents, proxyUnlockEvents := make([]*models.ProxyLockEvent, 0), make([]*models.ProxyUnlockEvent, 0)
	nftProxyLockEvents, nftProxyUnlockEvents, err := g.getNFTProxyEventByBlockNumber(g.platonCfg.NFTProxyContract, height, height)
	if err != nil {
		return nil, nil, nil, 0, 0, err
	}

	proxyLockEvents = append(proxyLockEvents, nftProxyLockEvents...)
	proxyUnlockEvents = append(proxyUnlockEvents, nftProxyUnlockEvents...)

	srcTransactions := make([]*models.SrcTransaction, 0)
	dstTransactions := make([]*models.DstTransaction, 0)
	for _, lockEvent := range eccmLockEvents {
		if lockEvent.Method == utils.Crosschainlock {
			//logs.Info("(lock) from chain: %s, txhash: %s, txid: %s", g.GetChainName(), lockEvent.TxHash, lockEvent.Txid)
			srcTransaction := &models.SrcTransaction{}
			srcTransaction.ChainID = g.GetChainID()
			srcTransaction.Hash = lockEvent.TxHash
			srcTransaction.State = 1
			srcTransaction.Fee = models.NewBigIntFromInt(int64(lockEvent.Fee))
			srcTransaction.Time = time
			srcTransaction.Height = lockEvent.Height
			srcTransaction.User = lockEvent.User
			srcTransaction.DstChainID = uint64(lockEvent.Tchain)
			srcTransaction.Contract = lockEvent.Contract
			srcTransaction.Key = lockEvent.Txid
			srcTransaction.Param = hex.EncodeToString(lockEvent.Value)
			var lock *models.ProxyLockEvent
			for _, v := range proxyLockEvents {
				if v.TxHash == lockEvent.TxHash {
					lock = v
					break
				}
			}
			if lock != nil {
				toAssetHash := lock.ToAssetHash
				srcTransfer := &models.SrcTransfer{}
				srcTransfer.Time = time
				srcTransfer.ChainID = g.GetChainID()
				srcTransfer.TxHash = lockEvent.TxHash
				srcTransfer.From = lockEvent.User
				srcTransfer.To = lockEvent.Contract
				srcTransfer.Asset = lock.FromAssetHash
				srcTransfer.TokenID = models.NewBigInt(lock.TokenID)
				srcTransfer.DstChainID = uint64(lock.ToChainID)
				srcTransfer.DstAsset = toAssetHash
				srcTransfer.DstUser = lock.ToAddress
				srcTransaction.SrcTransfer = srcTransfer
				if g.isNFTECCMLockEvent(lockEvent) {
					srcTransaction.Standard = models.TokenTypeErc721
					srcTransaction.SrcTransfer.Standard = models.TokenTypeErc721
				}
			}
			if srcTransaction.SrcTransfer != nil || srcTransaction.SrcSwap != nil {
				srcTransactions = append(srcTransactions, srcTransaction)
			}
		}
	}
	// save unLockEvent to db
	for _, unLockEvent := range eccmUnLockEvents {
		if unLockEvent.Method == utils.Crosschainunlock {
			//logs.Info("(unlock) to chain: %s, txhash: %s", g.GetChainName(), unLockEvent.TxHash)
			dstTransaction := &models.DstTransaction{}
			dstTransaction.ChainID = g.GetChainID()
			dstTransaction.Hash = unLockEvent.TxHash
			dstTransaction.State = 1
			dstTransaction.Fee = models.NewBigIntFromInt(int64(unLockEvent.Fee))
			dstTransaction.Time = time
			dstTransaction.Height = unLockEvent.Height
			dstTransaction.SrcChainID = uint64(unLockEvent.FChainID)
			dstTransaction.Contract = unLockEvent.Contract
			dstTransaction.PolyHash = unLockEvent.RTxHash
			var unlock *models.ProxyUnlockEvent
			for _, v := range proxyUnlockEvents {
				if v.TxHash == unLockEvent.TxHash {
					unlock = v
					break
				}
			}
			if unlock != nil {
				dstTransfer := &models.DstTransfer{}
				dstTransfer.TxHash = unLockEvent.TxHash
				dstTransfer.Time = time
				dstTransfer.ChainID = g.GetChainID()
				dstTransfer.From = unLockEvent.Contract
				dstTransfer.To = unlock.ToAddress
				dstTransfer.Asset = unlock.ToAssetHash
				dstTransfer.TokenID = models.NewBigInt(unlock.TokenID)
				dstTransaction.DstTransfer = dstTransfer
				if g.isNFTECCMUnlockEvent(unLockEvent) {
					dstTransaction.Standard = models.TokenTypeErc721
					dstTransaction.DstTransfer.Standard = models.TokenTypeErc721
				}
			}
			if dstTransaction.DstTransfer != nil {
				dstTransactions = append(dstTransactions, dstTransaction)
			}
		}
	}
	return wrapperTransactions, srcTransactions, dstTransactions, len(proxyLockEvents), len(proxyUnlockEvents), nil
}

func (g *PlatonChainListen) getNFTWrapperEventByBlockNumber(wrapAddrStr string, startHeight, endHeight uint64) ([]*models.WrapperTransaction, error) {
	wrapAddr := common.HexToAddress(wrapAddrStr)
	wrapperContract, err := nftwrap.NewPolyNFTWrapper(wrapAddr, g.platonSdk.GetClient())
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractEventByBlock, error: %s", err.Error())
	}
	opt := &bind.FilterOpts{
		Start:   startHeight,
		End:     &endHeight,
		Context: context.Background(),
	}

	// get geth lock events from given block
	wrapperTransactions := make([]*models.WrapperTransaction, 0)
	lockEvents, err := wrapperContract.FilterPolyWrapperLock(opt, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractEventByBlock, filter lock events :%s", err.Error())
	}
	for lockEvents.Next() {
		evt := lockEvents.Event
		wtx := utils.WrapLockEvent2WrapTx(evt)
		wrapperTransactions = append(wrapperTransactions, wtx)
	}
	speedupEvents, err := wrapperContract.FilterPolyWrapperSpeedUp(opt, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractEventByBlock, filter lock events :%s", err.Error())
	}
	for speedupEvents.Next() {
		evt := speedupEvents.Event
		wtx := utils.WrapSpeedUpEvent2WrapTx(evt)
		wrapperTransactions = append(wrapperTransactions, wtx)
	}
	return wrapperTransactions, nil
}

func (g *PlatonChainListen) GetConsumeGas(hash common.Hash) uint64 {
	tx, err := g.platonSdk.GetTransactionByHash(hash)
	if err != nil {
		return 0
	}
	receipt, err := g.platonSdk.GetTransactionReceipt(hash)
	if err != nil {
		return 0
	}
	return tx.GasPrice().Uint64() * receipt.GasUsed
}

func (g *PlatonChainListen) getECCMEventByBlockNumber(contractAddr string, startHeight uint64, endHeight uint64) ([]*models.ECCMLockEvent, []*models.ECCMUnlockEvent, error) {
	eccmContractAddress := common.HexToAddress(contractAddr)
	eccmContract, err := eccm.NewEthCrossChainManager(eccmContractAddress, g.platonSdk.GetClient())
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, error: %s", err.Error())
	}
	opt := &bind.FilterOpts{
		Start:   startHeight,
		End:     &endHeight,
		Context: context.Background(),
	}
	// get ethereum lock events from given block
	eccmLockEvents := make([]*models.ECCMLockEvent, 0)
	crossChainEvents, err := eccmContract.FilterCrossChainEvent(opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, filter lock events :%s", err.Error())
	}
	for crossChainEvents.Next() {
		evt := crossChainEvents.Event
		Fee := g.GetConsumeGas(evt.Raw.TxHash)
		eccmLockEvents = append(eccmLockEvents, &models.ECCMLockEvent{
			Method:   utils.Crosschainlock,
			Txid:     hex.EncodeToString(evt.TxId),
			TxHash:   evt.Raw.TxHash.String()[2:],
			User:     strings.ToLower(evt.Sender.String()[2:]),
			Tchain:   uint32(evt.ToChainId),
			Contract: strings.ToLower(evt.ProxyOrAssetContract.String()[2:]),
			Value:    evt.Rawdata,
			Height:   evt.Raw.BlockNumber,
			Fee:      Fee,
		})
	}
	// ethereum unlock events from given block
	eccmUnlockEvents := make([]*models.ECCMUnlockEvent, 0)
	executeTxEvent, err := eccmContract.FilterVerifyHeaderAndExecuteTxEvent(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, filter unlock events :%s", err.Error())
	}

	for executeTxEvent.Next() {
		evt := executeTxEvent.Event
		Fee := g.GetConsumeGas(evt.Raw.TxHash)
		eccmUnlockEvents = append(eccmUnlockEvents, &models.ECCMUnlockEvent{
			Method: utils.Crosschainunlock,
			TxHash: evt.Raw.TxHash.String()[2:],
			//RTxHash:  utils.HexStringReverse(hex.EncodeToString(evt.CrossChainTxHash)),
			Contract: hex.EncodeToString(evt.ToContract),
			FChainID: uint32(evt.FromChainID),
			Height:   evt.Raw.BlockNumber,
			Fee:      Fee,
		})
	}
	return eccmLockEvents, eccmUnlockEvents, nil
}

func (g *PlatonChainListen) getNFTProxyEventByBlockNumber(
	proxyAddrStr string,
	startHeight, endHeight uint64) (
	[]*models.ProxyLockEvent,
	[]*models.ProxyUnlockEvent,
	error,
) {
	proxyAddr := common.HexToAddress(proxyAddrStr)
	proxyContract, err := nftlp.NewPolyNFTLockProxy(proxyAddr, g.platonSdk.GetClient())
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, error: %s", err.Error())
	}
	opt := &bind.FilterOpts{
		Start:   startHeight,
		End:     &endHeight,
		Context: context.Background(),
	}
	// get ethereum lock events from given block
	proxyLockEvents := make([]*models.ProxyLockEvent, 0)
	lockEvents, err := proxyContract.FilterLockEvent(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, filter lock events :%s", err.Error())
	}
	for lockEvents.Next() {
		proxyLockEvent := utils.ConvertLockProxyEvent(lockEvents.Event)
		proxyLockEvents = append(proxyLockEvents, proxyLockEvent)
	}

	// ethereum unlock events from given block
	proxyUnlockEvents := make([]*models.ProxyUnlockEvent, 0)
	unlockEvents, err := proxyContract.FilterUnlockEvent(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, filter unlock events :%s", err.Error())
	}
	for unlockEvents.Next() {
		proxyUnlockEvent := utils.ConvertUnlockProxyEvent(unlockEvents.Event)
		proxyUnlockEvents = append(proxyUnlockEvents, proxyUnlockEvent)
	}
	return proxyLockEvents, proxyUnlockEvents, nil
}

func (g *PlatonChainListen) isNFTECCMLockEvent(event *models.ECCMLockEvent) bool {
	addr1 := common.HexToAddress(event.Contract)
	addr2 := common.HexToAddress(g.platonCfg.NFTProxyContract)
	return bytes.Equal(addr1.Bytes(), addr2.Bytes())
}

func (g *PlatonChainListen) isNFTECCMUnlockEvent(event *models.ECCMUnlockEvent) bool {
	addr1 := common.HexToAddress(event.Contract)
	addr2 := common.HexToAddress(g.platonCfg.NFTProxyContract)
	return bytes.Equal(addr1.Bytes(), addr2.Bytes())
}
