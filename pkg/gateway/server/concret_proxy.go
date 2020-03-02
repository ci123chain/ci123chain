package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"net/http"
)

type ConcretProxy struct {
	ProxyType types.ProxyType
	ResponseChannel chan []byte
}

func NewConcretProxy(pt types.ProxyType) *ConcretProxy {

	cp := &ConcretProxy{
		ProxyType: pt,
		ResponseChannel:make(chan []byte),
	}
	return cp
}


func (cp *ConcretProxy) Handle(r *http.Request, backends []types.Instance, RequestParams map[string]string) []byte {

	backendsLen := len(backends)
	var resultResp []Response
	for i := 0; i < backendsLen; i++ {
		fmt.Println(backends[i].URL().Host)
	}

	if backendsLen == 1 {
		resByte, _, err := SendRequest(backends[0].URL(), r, RequestParams)
		if err != nil {
			//
			err = errors.New("failed get response")
			res, _ := json.Marshal(types.ErrorResponse{
				Err:  err.Error(),
			})
			//return res
			cp.ResponseChannel <- res
			return res
		}
		cp.ResponseChannel <- resByte
		return resByte
	}else {
		for i := 0; i < backendsLen; i++ {
			var result Response
			resByte, _, err := SendRequest(backends[i].URL(),r, RequestParams)
			if err != nil {
				//return nil, errors.New("failed get response")
			}
			result = HandleResponse(resByte)
			resultResp = append(resultResp, result)
		}
	}
	resultByte, err := json.Marshal(resultResp)
	if err != nil {
		err = errors.New("failed to unmarshal response bytes")
		res, _ := json.Marshal(types.ErrorResponse{
			Err:  err.Error(),
		})
		cp.ResponseChannel <- res
		return res
	}

	cp.ResponseChannel <- resultByte
	return resultByte

/*
	//------
	if backendsLen == 1 {
		byte, _, err := SendRequest(backends[0].URL(), r)
		if err != nil {
			//
		}
		result := FormateResponse(byte)
		return result
	}

	clasterTaskPool := NewClasterTaskPool(3)
	clasterTaskPool.Run()
	go func() {
		for i := 0; i < backendsLen; i++ {
			//fmt.Println(i)
			job := NewClasterTask(backends[i].URL(), r, backendsLen, i)
			clasterTaskPool.JobQueue <- job
		}
	}()

	<- response
	allResult, err := json.Marshal(ConcretResultResponse)
	if err != nil {
		//
	}
	clasterTaskPool.Stop()
	return allResult
	*/
}

func (cp *ConcretProxy) Response() *chan []byte {
	return &cp.ResponseChannel
}

