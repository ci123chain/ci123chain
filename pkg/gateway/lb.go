package gateway

import (
	"flag"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/gateway/backend"
	"github.com/tanhuiya/ci123chain/pkg/gateway/couchdbsource"
	"github.com/tanhuiya/ci123chain/pkg/gateway/lbpolicy"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Attempts int = iota
	Retry
)


// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

// GetAttemptsFromContext returns the attempts for request
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

// lb load balances the incoming request
func lb(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := serverPool.GetNextPeer()
	if peer != nil {
		peer.Proxy().ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// isAlive checks whether a backend is Alive by establishing a TCP connection
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	_ = conn.Close()
	return true
}

// healthCheck runs a routine for check status of the backends every 2 mins
func healthCheck() {
	t := time.NewTicker(time.Second * 15)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

func fetchSharedRoutine()  {
	serverPool.SharedCheck()
	t := time.NewTicker(time.Second * 14)
	for {
		select {
		case <-t.C:
			serverPool.SharedCheck()
		}
	}
}


var serverPool *ServerPool

func Start() {
	var serverList string
	var statedb, dbname string
	var port int
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.StringVar(&statedb, "statedb", "couchdb://couchdb-service:5984", "server resource")
	flag.StringVar(&dbname, "db", "ci123", "db name")

	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()
	//
	//if len(serverList) == 0 {
	//	log.Fatal("Please provide one or more backends to load balance")
	//}
	policy := lbpolicy.NewRoundPolicy()
	svr := couchdbsource.NewCouchSource(dbname, statedb)

	serverPool = NewServerPool(backend.NewBackEnd ,lb, policy, svr)

	list := strings.Split(serverList, ",")
	serverPool.ConfigServerPool(list)
	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb),
	}

	// start health checking
	go healthCheck()

	go fetchSharedRoutine()

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
