package gateway

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/gateway/logger"
	"github.com/tanhuiya/ci123chain/pkg/gateway/server"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"net/http"
)

type SpecificJob struct {
	Request        *http.Request
	Proxy          types.Proxy
	Backends       []types.Instance
	RequestParams    map[string]string

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

	resultBytes := sjob.Proxy.Handle(sjob.Request, sjob.Backends, sjob.RequestParams)

	logger.Info("===\n Request for : %s; Params: %v;  response: %v", sjob.Request.URL.String(), sjob.RequestParams, string(resultBytes))
}

func NewSpecificJob(r *http.Request, backends []types.Instance) *SpecificJob {

	proxy, err, reqParams := ParseURL(r)
	if err != nil {
		return nil
	}

	job := &SpecificJob{
		Request: r,
		Proxy:   proxy,
		Backends:backends,
		RequestParams:reqParams,
	}
	job.ResponseChan = proxy.Response()

	return job
}

func ParseURL(r *http.Request) (types.Proxy, error, map[string]string){
	//body, _ := ioutil.ReadAll(r.Body)
	//data := r.Form

	/*
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
	*/
	//data := r.PostForm
	//data := url.Values{}
	r.ParseForm()
	var data = map[string]string{}

	for k, v := range r.Form {
		//fmt.Println("key is: ", k)
		//fmt.Println("val is: ", v)
		//data.Set(k, v[0])
		key := k
		value := v[0]
		data[key] = value
	}

	//data.Set("height", "3")
	params := r.FormValue("proxy")

	pt := types.ProxyType(params)
	switch pt {
	case types.LB:
		return server.NewLBProxy(pt), nil, data
	case types.Concret:
		return server.NewConcretProxy(pt), nil, data
	case types.Filter:
		return server.NewFilterProxy(pt), nil, data
	default:
		return server.NewLBProxy(pt), nil, data
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