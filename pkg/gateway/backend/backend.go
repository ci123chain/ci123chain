package backend

import (
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"net/url"
	"sync"
)

// Backend holds the data about a server
type Backend struct {
	url          *url.URL
	Alive        bool
	mux          sync.RWMutex
	retry 		 int
}

// SetAlive for this backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	if alive {
		b.retry = 0
	} else {
		b.retry++
	}
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}

func (b *Backend) URL() *url.URL{
	return b.url
}

func (b *Backend) FailTime() int {
	return b.retry
}

func NewBackEnd(url *url.URL, alive bool) types.Instance {
	return &Backend{
		url:          url,
		Alive:        alive,
		retry: 		  0,
	}
}

