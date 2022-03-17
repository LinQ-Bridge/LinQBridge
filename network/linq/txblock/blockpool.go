package txblock

import (
	"sync"

	"land-bridge/models"
	"land-bridge/network/bridge"
)

type TxInfo struct {
	W       *models.WrapperTransaction
	TxParam *bridge.TxParam
}

type BlockPool struct {
	mux sync.RWMutex
	m   map[string]*TxInfo
	l   []string
}

func NewBlockPool() *BlockPool {
	return &BlockPool{
		m: make(map[string]*TxInfo),
		l: []string{},
	}
}

func (p *BlockPool) Get(key string) (bool, *TxInfo) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	if v, ok := p.m[key]; ok {
		return true, v
	} else {
		return false, nil
	}
}

func (p *BlockPool) Push(key string, value *TxInfo) {
	p.mux.Lock()
	defer p.mux.Unlock()
	if _, ok := p.m[key]; !ok {
		p.l = append(p.l, key)
		p.m[key] = value
	}
}

func (p *BlockPool) First() *TxInfo {
	p.mux.RLock()
	defer p.mux.RUnlock()

	if len(p.l) > 0 {
		v, _ := p.m[p.l[0]]
		return v
	} else {
		return nil
	}
}

func (p *BlockPool) Delete(key string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if _, ok := p.m[key]; ok {
		delete(p.m, key)
		for i, s := range p.l {
			if s == key {
				p.l = append(p.l[:i], p.l[i+1:]...)
				return
			}
		}
	}
}

func (p *BlockPool) Clean() {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.m = make(map[string]*TxInfo)
	p.l = []string{}
}
