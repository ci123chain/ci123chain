package gateway

import (
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/gateway/server"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"io/ioutil"
	"net/http"
	"strings"
)

type SpecificJob struct {
	Request        *http.Request
	Proxy          types.Proxy
	ResponseWriter http.ResponseWriter
	Backends       []types.Instance
}

type OtherParams struct {
	ID   uint64     `json:"id"`
	Host string     `json:"host"`
}

type Params struct {
	Proxy string      `json:"proxy"`
	Other  OtherParams `json:"other"`
}

type RequestParams struct {
	Proxy string `json:"proxy"`
	Data interface{} `json:"data"`
}

type NewRequestParams struct {
	Data interface{} `json:"data"`
}


func (sjob *SpecificJob) Do() {
	resultBytes, err := sjob.Proxy.Handle(sjob.Request, sjob.Backends)
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

	proxy, err, newRequest := ParseURL(r)
	if err != nil {
		_, _ = w.Write([]byte("unexpected proxy"))
	}
	r = newRequest

	return &SpecificJob{
		Request: r,
		Proxy:   proxy,
		Backends:backends,
		ResponseWriter:w,
	}
}

func ParseURL(r *http.Request) (types.Proxy, error, *http.Request){
	body, _ := ioutil.ReadAll(r.Body)
	var params RequestParams
	err := json.Unmarshal(body, &params)
	if err != nil {
		return server.NewErrProxy("err"), err, nil
	}
	nrp := NewRequestParams{Data:params.Data}
	newByte, err := json.Marshal(nrp)
	if err != nil {
		//
	}
	newReq, err := http.NewRequest(r.Method, r.Host, strings.NewReader(string(newByte)))
	if err != nil {
		//
	}

	pt := types.ProxyType(params.Proxy)
	switch params.Proxy {
	case types.LB:
		return server.NewLBProxy(pt), nil, newReq
	case types.Concret:
		return server.NewConcretProxy(pt), nil, newReq
	case types.Filter:
		return server.NewFilterProxy(pt), nil, newReq
	default:
		return server.NewErrProxy("unexpected policy"), errors.New("unexpected policy"), nil
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