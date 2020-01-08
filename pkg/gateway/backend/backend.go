package backend


import (
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Backend holds the data about a server
type Backend struct {
	url          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// SetAlive for this backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
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

func (b *Backend) Proxy() *httputil.ReverseProxy {
	return b.ReverseProxy
}

func NewBackEnd(url *url.URL, alive bool, proxy *httputil.ReverseProxy) types.Instance {
	return &Backend{
		url:          url,
		Alive:        alive,
		ReverseProxy: proxy,
	}
}