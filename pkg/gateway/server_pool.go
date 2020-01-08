package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

const Name  = "ci123"
const SharedKey  = "order//OrderBook"
//const StateDB = "couchdb://couchdb-service:5984"
const StateDB = "couchdb://192.168.2.89:30301"


func (s *ServerPool)ConfigServerPool(tokens []string)  {
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
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

		serverPool.AddBackend(&Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}
}

func NewServerPool(lb http.HandlerFunc) *ServerPool {
	return &ServerPool{
		backends: make([]*Backend, 0),
		current: 0,
		lb:lb,
	}
}

// ServerPool holds information about reachable backends
type ServerPool struct {
	backends []*Backend
	current  uint64
	lb 		 http.HandlerFunc
}

func (s *ServerPool) SharedCheck()  {
	conn, err := s.GetDBConnection()
	if err != nil {
		log.Println(err)
	}
	bz := conn.Get([]byte(SharedKey))
	var shared map[string]interface{}
	err = json.Unmarshal(bz, &shared)
	if err != nil {
		log.Println(err)
	}

	orderDict, ok := shared["value"].(map[string]interface{})
	if !ok {
		return
	}
	lists := orderDict["lists"].([]interface{})
	if !ok {
		return
	}
	var hostArr []string
	for _, value := range lists {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		name := item["name"].(string)
		hostArr = append(hostArr, name)
	}
	log.Println(hostArr)
	// todo
	s.ConfigServerPool(hostArr)
}

func (svr *ServerPool)GetDBConnection() (db *couchdb.GoCouchDB, err error) {
	s := strings.Split(StateDB, "://")
	if len(s) < 2 {
		return nil, errors.New("statedb format error")
	}
	if s[0] != "couchdb" {
		return nil, errors.New("statedb format error")
	}
	auths := strings.Split(s[1], "@")

	if len(auths) < 2 {
		db, err = couchdb.NewGoCouchDB(Name, auths[0],nil)
	} else {
		info := auths[0]
		userpass := strings.Split(info, ":")
		if len(userpass) < 2 {
			db, err = couchdb.NewGoCouchDB(Name, auths[1],nil)
		}
		auth := &couchdb.BasicAuth{Username: userpass[0], Password: userpass[1]}
		db, err = couchdb.NewGoCouchDB(Name, auths[1], auth)
	}
	return
}

// AddBackend to the server pool
func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

// NextIndex atomically increase the counter and return an index
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// MarkBackendStatus changes a status of a backend
func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() *Backend {
	// loop entire backends to find out an Alive backend
	next := s.NextIndex()
	l := len(s.backends) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(s.backends) // take an index by modding
		if s.backends[idx].IsAlive() { // if we have an alive backend, use it and store if its not the original one
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

// HealthCheck pings the backends and update the status
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}
