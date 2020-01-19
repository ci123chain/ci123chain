package gateway

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/gateway/logger"
	"github.com/tanhuiya/ci123chain/pkg/gateway/server"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"io/ioutil"
	"net/http"
)

type SpecificJob struct {
	Request        *http.Request
	Proxy          types.Proxy
	Backends       []types.Instance
	RequestBody    []byte

	ResponseChan 	*chan []byte
}


func (sjob *SpecificJob) Do() {
	if len(sjob.Backends) < 1 {
		res, _ := json.Marshal(types.ErrorResponse{
			Err:  "service backend not found",
		})
		*sjob.ResponseChan <- res
		return
	}

	resultBytes := sjob.Proxy.Handle(sjob.Request, sjob.Backends, sjob.RequestBody)

	logger.Info("===\n Request for : %s; Params: %v;  response: %v", sjob.Request.URL.String(), sjob.RequestBody, string(resultBytes))
}

func NewSpecificJob(r *http.Request, backends []types.Instance) *SpecificJob {

	proxy, err, reqBody := ParseURL(r)
	if err != nil {
		return nil
	}

	job := &SpecificJob{
		Request: r,
		Proxy:   proxy,
		Backends:backends,
		RequestBody:reqBody,
	}
	job.ResponseChan = proxy.Response()

	return job
}

func ParseURL(r *http.Request) (types.Proxy, error, []byte){
	body, _ := ioutil.ReadAll(r.Body)

	var params types.RequestParams
	err := json.Unmarshal(body, &params)
	if err != nil {
		return server.NewErrProxy("err"), err, nil
	}

	nrp := types.NewRequestParams{Data:params.Data}
	newByte, err := json.Marshal(nrp)
	if err != nil {
		return server.NewErrProxy("err"), err, nil
	}

	pt := types.ProxyType(params.Proxy)
	switch params.Proxy {
	case types.LB:
		return server.NewLBProxy(pt), nil, newByte
	case types.Concret:
		return server.NewConcretProxy(pt), nil, newByte
	case types.Filter:
		return server.NewFilterProxy(pt), nil, newByte
	default:
		return server.NewLBProxy(pt), nil, newByte
	}
}

type Worker struct {
	JobQueue chan types.Job
}

func NewWorker() Worker {
	return Worker{JobQueue: make(chan types.Job)}
}

func (w Worker) Run(wq chan chan types.Job) {
	go func() {
		for {
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				job.Do()
			}
		}
	}()
}