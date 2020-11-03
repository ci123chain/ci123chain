package server

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"net/http"
)

type FilterProxy struct {
	ProxyType types.ProxyType
	ResponseChannel chan []byte
}

func NewFilterProxy(pt types.ProxyType) *FilterProxy {

	fp :=  &FilterProxy{
		ProxyType:pt,
		ResponseChannel:make(chan []byte),
	}
	return fp
}

func (fp *FilterProxy) Handle(r *http.Request, backEnds []types.Instance, RequestParams map[string]string) []byte {

	backEndsLen := len(backEnds)
	var resultByte []byte
	var resultResp []Response
	var Responses []string
	var existResult = false

	if backEndsLen == 1 {
		resByte, _, err := SendRequest(backEnds[0].URL(), r, RequestParams)
		if err != nil {
			err = errors.New("failed get response")
			res, _ := json.Marshal(types.ErrorResponse{
				Ret: 0,
				Message:  err.Error(),
			})
			fp.ResponseChannel <- res
			return res
		}
		fp.ResponseChannel <- resByte
		return resByte
	}
	clasterTaskPool := NewClasterTaskPool(3)
	clasterTaskPool.Run()
	var task *ClasterTask
	var tasks []*ClasterTask
	for i := 0; i < backEndsLen; i++ {
		task = NewClasterTask(backEnds[i].URL(), r, RequestParams)
		clasterTaskPool.JobQueue <- task
		tasks = append(tasks, task)
	}

	for j := range tasks {
		select {
		case resByte := <- tasks[j].responseChan:
			Responses = append(Responses, string(resByte))
			result := HandleResponse(resByte)
			resultResp = append(resultResp, result)
		}
		if len(resultResp) == backEndsLen {
			break
		}
	}
	if len(Responses) == 0 || len(resultResp) == 0 {
		err := errors.New("responses is empty")
		res, _ := json.Marshal(types.ErrorResponse{
			Ret: 0,
			Message:  err.Error(),
		})
		fp.ResponseChannel <- res
		return res
	}

	for i := range resultResp {
		if resultResp[i].Message == "" && resultResp[i].Data != nil {
			resultByte, _ = json.Marshal(resultResp[i])
			existResult = true
			break
		}
	}
	if existResult == false {
		resultByte = []byte(Responses[0])
	}
	fp.ResponseChannel <- resultByte
	return resultByte

}

func (fp *FilterProxy) Response() *chan []byte {
	return &fp.ResponseChannel
}