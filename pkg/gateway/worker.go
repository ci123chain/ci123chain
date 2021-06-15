package gateway

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/gateway/server"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"io/ioutil"
	"net/http"
	"reflect"
)

type SpecificJob struct {
	Request        *http.Request
	Proxy          types.Proxy
	BackEnds       []types.Instance
	RequestParams    map[string]string

	ResponseChan 	*chan []byte
}


func (sjob *SpecificJob) Do() {
	if reflect.TypeOf(sjob.Proxy).Elem() != reflect.TypeOf(server.DeployProxy{}) {
		if len(sjob.BackEnds) < 1 {
			res, _ := json.Marshal(types.ErrorResponse{
				Ret: 0,
				Message:  "service backend not found",
			})
			*sjob.ResponseChan <- res
			return
		}
	}
	// for debug
	//resultBytes := sjob.Proxy.Handle(sjob.Request, sjob.BackEnds, sjob.RequestParams)
	//logger.Info("===\n Request for : %s; Params: %v;  response: %v", sjob.Request.URL.String(), sjob.RequestParams, string(resultBytes))
}

func NewSpecificJob(r *http.Request, backends []types.Instance) *SpecificJob {

	proxy, err, reqParams := ParseURL(r)
	if err != nil {
		return nil
	}
	job := &SpecificJob{
		Request: r,
		Proxy:   proxy,
		BackEnds:backends,
		RequestParams:reqParams,
	}
	job.ResponseChan = proxy.Response()

	return job
}

func ParseURL(r *http.Request) (types.Proxy, error, map[string]string){
	var data = map[string]string{}
	var codeStr string
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		for k, v := range r.PostForm {
			key := k
			value := v[0]
			data[key] = value
		}
	}else {
		if r.MultipartForm != nil {
			for k, v := range r.MultipartForm.Value {
				key := k
				value := v[0]
				data[key] = value
			}
			file, _, err := r.FormFile("code_file")
			if err != nil {
				codeStr = ""
			} else {
				wasmcode, err := ioutil.ReadAll(file)
				if err != nil {
					codeStr = ""
				}else {
					codeStr = hex.EncodeToString(wasmcode)
				}
			}
			data["code_str"] = codeStr
		}
	}
	params := r.FormValue("proxy")

	pt := types.ProxyType(params)
	switch pt {
	case types.LB:
		return server.NewLBProxy(pt), nil, data
	case types.Concret:
		return server.NewConcretProxy(pt), nil, data
	case types.Filter:
		return server.NewFilterProxy(pt), nil, data
	case types.Deploy:
		return server.NewDeployProxy(pt), nil, data
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