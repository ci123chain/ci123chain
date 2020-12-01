package gateway

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"net/url"
)


type BackendProto func(url *url.URL, alive bool) types.Instance

type ServerPool struct{
	backendProto 	BackendProto
	backends 		[]types.Instance
	svrsource 		types.ServerSource
	workerlen       int
	JobQueue        chan types.Job
	WorkerQueue     chan chan types.Job
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

func (s *ServerPool) ConfigServerPool(tokens []string)  {
	for _, tok := range tokens {
		exist := false
		for _, back := range s.backends {
			if back.URL().String() == tok {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		serverUrl, err := url.Parse(tok)
		if err != nil {
			logger.Info("info %s", err)
		}
		
		if !isBackendAlive(serverUrl) {
			continue
		}
/*
		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			logger.Debug("[%s] %s\n", serverUrl.Host, e.Error())
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
			logger.Debug("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			AllHandle(writer, request.WithContext(ctx))
		}
		*/

		serverPool.AddBackend(s.backendProto(serverUrl, true))
		logger.Debug("Configured server: %s\n", serverUrl)
	}
}

func NewServerPool(backProto BackendProto, svrsource types.ServerSource, workerlen int) *ServerPool {

	return &ServerPool{
		backendProto: 	backProto,
		backends: 		make([]types.Instance, 0),
		svrsource: 		svrsource,
		workerlen:      workerlen,
		JobQueue:       make(chan types.Job),
		WorkerQueue:    make(chan chan types.Job, workerlen),
	}
}

func (s *ServerPool) Run() {
	for i := 0; i < s.workerlen; i++ {

		worker := NewWorker()
		worker.Run(s.WorkerQueue)
	}
	// get new job, send to worker.
	go func() {
		for {
			select {
			case job := <-s.JobQueue:
				worker := <-s.WorkerQueue
				worker <- job
			}
		}
	}()
}

func (s *ServerPool) SharedCheck() {
	hosts := s.svrsource.FetchSource()
	if len(hosts) > 0 {
		s.ConfigServerPool(hosts)
	}
}

func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL())
		b.SetAlive(alive)
		if !alive {
			status = "down"
			logger.Warn("%s [%s]\n", b.URL(), status)
		}
		logger.Info("%s [%s]\n", b.URL(), status)
	}
}

func (s *ServerPool) getDeadList() []string{
	if len(s.backends) == 0 {
		return nil
	}

	var deadList []string
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL())
		b.SetAlive(alive)
		if !alive {
			status = "down"
			deadList = append(deadList, fmt.Sprintf("%s [%s]", b.URL(), status))
			logger.Warn("%s [%s]\n", b.URL(), status)
		}
		logger.Info("%s [%s]\n", b.URL(), status)
	}
	return deadList
}

