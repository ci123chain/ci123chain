package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Response struct {
	Ret 	interface{} 	`json:"ret"`
	Data 	interface{}	`json:"data"`
	Message	string	`json:"message"`
}

func SendRequest(url *url.URL,r *http.Request, reqBody []byte) ([]byte, *http.Response, error) {

	cli := &http.Client{}
	//body := make([]byte, 0)
	reqUrl := "http://" + url.Host + r.URL.Path


	req2, err := http.NewRequest(r.Method, reqUrl, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, nil, err
	}

	// set request content type
	contentType := r.Header.Get("Content-Type")
	req2.Header.Set("Content-Type", contentType)
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
		//
	}
	return res
}




var sc = make(chan int)

type ClasterJob interface {
	Do()
}

type ClasterTask struct {
	num int
	id int
	url *url.URL
	r *http.Request
}

func NewClasterTask(url *url.URL, r *http.Request, num, id int) *ClasterTask{
	return &ClasterTask{
		url:url,
		r:r,
		num:num,
		id:id,
	}
}


func (ct *ClasterTask) Do() {
	//
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
	go w.stop()
}

func (w Worker)stop() {
	<- sc
	w.JobQueue = nil
}


type ClasterTaskPool struct {

	workerlen   int
	JobQueue    chan ClasterJob
	WorkerQueue chan chan ClasterJob
}

func NewClasterTaskPool(workerlen int) *ClasterTaskPool{

	return &ClasterTaskPool{
		workerlen:   workerlen,
		JobQueue:    make(chan ClasterJob),
		WorkerQueue: make(chan chan ClasterJob, workerlen),
	}
}

func (ctp *ClasterTaskPool) Run() {
	for i := 0; i < ctp.workerlen; i++ {
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

func (ctp *ClasterTaskPool) Stop() {
	sc <- 0
	ctp.JobQueue = nil
	ctp.WorkerQueue = nil
}