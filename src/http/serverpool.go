package http

import (
	"github.com/shaojintian/load_balancer/src/core"
	"log"
	"sync/atomic"
)

type ServerPool struct {
	backends []*core.Backend
	// current alive peer
	current uint64
}

func (sp *ServerPool) NextIndex() int {
	// % len == [0,len-1] ==  索引backends  == 约束在backends范围内
	if len(sp.backends) <= 0 {
		log.Fatalln("no backend in backends ")
	}
	return int(atomic.AddUint64(&sp.current, uint64(1))) % len(sp.backends)
}

func (sp *ServerPool) GetNextPeer() *core.Backend {
	nextIndex := sp.NextIndex()
	// round robin a full cycle
	cycleLen := len(sp.backends) + nextIndex
	for i := nextIndex; i < cycleLen; i++ {
		index := i % len(sp.backends)
		//Alive peer
		if sp.backends[index].IsAlive() {
			// spin current alive peer
			atomic.StoreUint64(&sp.current, uint64(index))
			return sp.backends[index]
		}

	}

	// all dead peers
	log.Println("all dead peers")
	return nil

}
