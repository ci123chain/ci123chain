package lbpolicy

import (
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"sync/atomic"
)

type RoundPolicy struct {
	current uint64    `json:"current"`
}

func NewRoundPolicy() *RoundPolicy {
	return &RoundPolicy{current:0}
}

func (rp *RoundPolicy) NextPeer(backends []types.Instance) types.Instance {
	next := rp.NextIndex(len(backends))
	l := len(backends) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(backends) // take an index by modding
		if backends[idx].IsAlive() { // if we have an alive backend, use it and store if its not the original one
			if i != next {
				atomic.StoreUint64(&rp.current, uint64(idx))
			}
			return backends[idx]
		}
	}
	return nil
}

func (rp *RoundPolicy) Current() int {
	return int(rp.current)
}

func (rp *RoundPolicy) NextIndex(round int) int {
	return int(atomic.AddUint64(&rp.current, uint64(1)) % uint64(round))
}