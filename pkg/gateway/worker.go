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
}

type OtherParams struct {
	ID   uint64     `json:"id"`
	Host string     `json:"host"`
}

type Params struct {
	Policy string      `json:"policy"`
	Other  OtherParams `json:"other"`
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
		sjob.ResponseWriter.Write(res)
	}else {
		sjob.ResponseWriter.Write(resultBytes)
	}
}

func NewSpecificJob(w http.ResponseWriter, r *http.Request, backends []types.Instance) *SpecificJob {

	proxy, err := ParseURL(r)
	if err != nil {
		w.Write([]byte("unexpected policy"))
	}

	return &SpecificJob{
		Request: r,
		Proxy:   proxy,
		Backends:backends,
		ResponseWriter:w,
	}
}

func ParseURL(r *http.Request) (types.Proxy, error){
	body, _ := ioutil.ReadAll(r.Body)
	var params Params
	err := json.Unmarshal(body, &params)
	if err != nil {
		return server.NewErrProxy("err"), err
	}
	pt := types.ProxyType(params.Policy)
	switch params.Policy {
	case types.LB:
		return server.NewLBProxy(pt), nil
	case types.Concret:
		return server.NewConcretProxy(pt), nil
	case types.Filter:
		return server.NewFilterProxy(pt), nil
	default:
		return server.NewErrProxy("unexpected policy"), errors.New("unexpected policy")
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