package server

import (
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
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

func (fp *FilterProxy) Handle(r *http.Request, backends []types.Instance, RequestParams map[string]string) []byte {

	backendsLen := len(backends)
	var resultByte []byte
	var resultResp []Response
	var Responses []string
	var existResult = false

	if backendsLen == 1 {
		resByte, _, err := SendRequest(backends[0].URL(), r, RequestParams)
		if err != nil {
			//
			err = errors.New("failed to get response")
			res, _ := json.Marshal(types.ErrorResponse{
				Err:  err.Error(),
			})
			fp.ResponseChannel <- res
			return res
		}
		fp.ResponseChannel <- resByte
		return resByte

	}else {
		for i := 0; i < backendsLen; i++ {
			var result Response
			resByte, _, _ := SendRequest(backends[i].URL(),r, RequestParams)
			Responses = append(Responses, string(resByte))
			result = HandleResponse(resByte)
			resultResp = append(resultResp, result)
		}
	}
	for i := range resultResp {
		if resultResp[i].Message == "" && resultResp[i].Data != nil {
			resultByte, _ = json.Marshal(resultResp[i])
			existResult = true
			break
		}
	}
	if existResult == false {
		resultByte = []byte(Responses[len(Responses) - 1])
	}
	fp.ResponseChannel <- resultByte
	return resultByte

	/*
	clasterTaskPool := NewClasterTaskPool(3)
	clasterTaskPool.Run()
	go func() {
		for i := 0; i < backendsLen; i++ {
			fmt.Println(i)
			job := NewClasterTask(backends[i].URL(), r, backendsLen, i)
			clasterTaskPool.JobQueue <- job
		}
	}()

	<- response
	fmt.Println(ResultResponse)
	byte := filter(ResultResponse, false)
	//w.Write(byte)

	clasterTaskPool.Stop()
	return byte
	*/
}

func (fp *FilterProxy) Response() *chan []byte {
	return &fp.ResponseChannel
}