package gateway

import (
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/gateway/server"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"io/ioutil"
	"net/http"
)

type SpecificJob struct {
	Request        *http.Request
	Proxy          types.Proxy
	ResponseWriter http.ResponseWriter
	Backends       []types.Instance
	RequestBody    []byte
}


func (sjob *SpecificJob) Do() {
	resultBytes, err := sjob.Proxy.Handle(sjob.Request, sjob.Backends, sjob.RequestBody)
	sjob.ResponseWriter.Header().Set("Content-Type", "application/json")
	if err != nil {
		errRes := types.ErrorResponse{
			Err:  err.Error(),
			//Code: types.ErrGetErrorResponse,
		}
		res, _ := json.Marshal(errRes)
		_, _ = sjob.ResponseWriter.Write(res)
	}else {
		_, _ = sjob.ResponseWriter.Write(resultBytes)
	}
}

func NewSpecificJob(w http.ResponseWriter, r *http.Request, backends []types.Instance) *SpecificJob {

	proxy, err, reqBody := ParseURL(r)
	if err != nil {
		_, _ = w.Write([]byte("unexpected proxy"))
	}
	//r = newRequest

	return &SpecificJob{
		Request: r,
		Proxy:   proxy,
		Backends:backends,
		ResponseWriter:w,
		RequestBody:reqBody,
	}
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
		//
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
		return server.NewErrProxy("unexpected policy"), errors.New("unexpected policy"), newByte
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