package gateway

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
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

// isAlive checks whether a backend is Alive by establishing a TCP connection
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println(fmt.Sprintf("Site unreachable for host: %s, error: %v", u.String(), err))
		return false
	}
	_ = conn.Close()
	return true
}

// healthCheck runs a routine for check status of the backends every 2 mins
func healthCheck() {
	t := time.NewTicker(time.Second * 20)
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
	t := time.NewTicker(time.Second * 15)
	for {
		select {
		case <-t.C:
			serverPool.SharedCheck()
		}
	}
}

var serverPool *ServerPool

/*func Start() {
	var serverList string
	var statedb, dbname string
	var port int
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.StringVar(&statedb, "statedb", "couchdb://couchdb-service:5984", "server resource")
	flag.StringVar(&dbname, "db", "ci123", "db name")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()

	policy := lbpolicy.NewRoundPolicy()
	svr := couchdbsource.NewCouchSource(dbname, statedb)

	serverPool = NewServerPool(backend.NewBackEnd ,lb, policy, svr, true, 3)

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
}*/
