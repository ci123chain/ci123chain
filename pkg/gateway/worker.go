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
	//ResponseWriter ResWriter
}

type OtherParams struct {
	ID   uint64     `json:"id"`
	Host string     `json:"host"`
}

type Params struct {
	Policy string      `json:"policy"`
	Other  OtherParams `json:"other"`
}


type ResWriter func(w http.ResponseWriter, result []byte)

func (sjob SpecificJob) Do() {
	sjob.Proxy.Handle(sjob.ResponseWriter,sjob.Request, sjob.Backends)
	//fmt.Println(string(byte))
	//sjob.ResponseWriter.WriteHeader()
	//sjob.ResponseWriter.Write([]byte("bbb"))
}

func NewSpecificJob(w http.ResponseWriter, r *http.Request, backends []types.Instance) SpecificJob {
	//

	proxy, err := ParseURL(r)
	if err != nil {
		w.Write([]byte("unexpected policy"))
	}

	return SpecificJob{
		Request: r,
		Proxy:   proxy,
		Backends:backends,
		ResponseWriter:w,
	}
}

func ParseURL(r *http.Request) (types.Proxy, error){
	body, _ := ioutil.ReadAll(r.Body)

	//
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