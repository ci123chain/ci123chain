package gateway

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
)

type BackendProto func(url *url.URL, alive bool, proxy *httputil.ReverseProxy) types.Instance

func (s *ServerPool)ConfigServerPool(tokens []string)  {
	for _, tok := range tokens {
LOOP:
		for _, back := range s.backends {
			if back.URL().String() == tok {
				goto LOOP
			}
		}

		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		if !isBackendAlive(serverUrl) {
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := GetRetryFromContext(request)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			// after 3 retries, mark this backend as down
			serverPool.MarkBackendStatus(serverUrl, false)

			// if the same request routing for few attempts with different backends, increase the count
			attempts := GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			lb(writer, request.WithContext(ctx))
		}

		serverPool.AddBackend(s.backendProto(serverUrl, true, proxy))
		log.Printf("Configured server: %s\n", serverUrl)
	}
}

func NewServerPool(backProto BackendProto, lb http.HandlerFunc, police types.LBPolicy, svrsource types.ServerSource) *ServerPool {
	return &ServerPool{
		backendProto: 	backProto,
		backends: 		make([]types.Instance, 0),
		lb:				lb,
		policy:  		police,
		svrsource: 		svrsource,
	}
}

// ServerPool holds information about reachable backends
type ServerPool struct {
	backendProto 	BackendProto
	backends 		[]types.Instance
	lb 		  		http.HandlerFunc
	policy 	  		types.LBPolicy
	svrsource 		types.ServerSource
}


func (s *ServerPool) SharedCheck() {
	hosts := s.svrsource.FetchSource()
	if len(hosts) > 0 {
		s.ConfigServerPool(hosts)
	}
}

// AddBackend to the server pool
func (s *ServerPool) AddBackend(backend types.Instance) {
	s.backends = append(s.backends, backend)
}


// MarkBackendStatus changes a status of a backend
func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL().String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() types.Instance {
	return s.policy.NextPeer(s.backends)
}

// HealthCheck pings the backends and update the status
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL())
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL(), status)
	}
}
