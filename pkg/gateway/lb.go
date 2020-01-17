package gateway

import (
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/gateway/logger"
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
		logger.Warn(fmt.Sprintf("Site unreachable for host: %s, error: %v", u.String(), err))
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
			logger.Debug("Starting health check...")
			serverPool.HealthCheck()
			logger.Debug("Health check completed")
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
