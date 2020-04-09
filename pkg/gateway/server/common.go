package server

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	Ret 	interface{} 	`json:"ret"`
	Data 	interface{}	    `json:"data"`
	Message	string	        `json:"message"`
}

func SendRequest(requestUrl *url.URL,r *http.Request, RequestParams map[string]string) ([]byte, *http.Response, error) {

	cli := &http.Client{
		Transport:&http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*60)    //设置建立连接超时
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
	}
	reqUrl := "http://" + requestUrl.Host + r.URL.Path
	data := url.Values{}
	for k, v := range RequestParams {
		data.Set(k, v)
	}


	req2, err := http.NewRequest(r.Method, reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req2.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content type
	//contentType := r.Header.Get("Content-Type")
	//req2.Header.Set("Content-Type", contentType)
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request
	rep2, err := cli.Do(req2)
	if err != nil {
		return nil, nil, err
	}
	b, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		panic(err)
	}
	return b, rep2, nil
}


func HandleResponse(b []byte) Response {
	var res Response
	err := json.Unmarshal(b, &res)
	if err != nil {
		var res = Response{
			Ret:     nil,
			Data:    nil,
			Message: "failed to unmarshal response",
		}
		return res
	}
	return res
}


type ClasterJob interface {
	Do()
}

type ClasterTask struct {
	url *url.URL
	r *http.Request
	requestParmas map[string]string

	responseChan chan []byte
}

func NewClasterTask(url *url.URL, r *http.Request, requestParams map[string]string) *ClasterTask{
	return &ClasterTask{
		url:url,
		r:r,
		requestParmas:requestParams,
		responseChan:make(chan []byte),
	}
}


func (ct *ClasterTask) Do() {
	res, _, err := SendRequest(ct.url, ct.r, ct.requestParmas)
	if err != nil {
		res, _ := json.Marshal(types.ErrorResponse{
			Err:  err.Error(),
		})
		ct.responseChan <- res
		return
	}
	ct.responseChan <- res
}

type Worker struct {
	JobQueue chan ClasterJob
	StopChannel chan int
}

func NewWorker() Worker {
	return Worker{
		JobQueue:make(chan ClasterJob),
		StopChannel:make(chan int),
	}
}

func (w Worker) Run(ctp chan chan ClasterJob) {
	go func() {
		for {
			ctp <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				job.Do()
			}
		}
	}()
}


type ClasterTaskPool struct {

	workerLen   int
	JobQueue    chan ClasterJob
	WorkerQueue chan chan ClasterJob
}

func NewClasterTaskPool(workerlen int) *ClasterTaskPool{

	return &ClasterTaskPool{
		workerLen:   workerlen,
		JobQueue:    make(chan ClasterJob),
		WorkerQueue: make(chan chan ClasterJob, workerlen),
	}
}

func (ctp *ClasterTaskPool) Run() {
	for i := 0; i < ctp.workerLen; i++ {
		worker := NewWorker()
		worker.Run(ctp.WorkerQueue)
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			case job := <-ctp.JobQueue:
				worker := <-ctp.WorkerQueue
				worker <- job
			}
		}
	}()
}