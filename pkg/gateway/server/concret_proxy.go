package server

import (
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"net/http"
)

type ConcretProxy struct {
	ProxyType types.ProxyType
}

func NewConcretProxy(pt types.ProxyType) *ConcretProxy {

	cp := &ConcretProxy{
		ProxyType: pt,
	}
	return cp
}


func (cp *ConcretProxy) Handle(r *http.Request, backends []types.Instance, reqBody []byte) ([]byte, error) {

	backendsLen := len(backends)
	var resultResp []ciRes

	if backendsLen == 1 {
		resByte, _, err := SendRequest(backends[0].URL(), r, reqBody)
		if err != nil {
			//
			return nil, errors.New("failed get response")
		}
		return resByte, nil
	}else {
		for i := 0; i < backendsLen - 1; i++ {
			var result ciRes
			resByte, _, err := SendRequest(backends[i].URL(),r, reqBody)
			if err != nil {
				//return nil, errors.New("failed get response")
			}
			result = AddResponses(resByte)
			resultResp = append(resultResp, result)
		}
	}
	resultByte, err := json.Marshal(resultResp)
	if err != nil {
		return nil, errors.New("failed to unmarshal response bytes")
	}
	return resultByte, nil

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

