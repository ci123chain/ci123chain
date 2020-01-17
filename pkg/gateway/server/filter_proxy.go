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

func (fp *FilterProxy) Handle(r *http.Request, backends []types.Instance, reqBody []byte) {

	backendsLen := len(backends)
	var resultResp []*http.Response
	var resultByte []byte
	var result []string

	if backendsLen == 1 {
		resByte, _, err := SendRequest(backends[0].URL(), r, reqBody)
		if err != nil {
			//
			err = errors.New("failed to get response")
			res, _ := json.Marshal(types.ErrorResponse{
				Err:  err.Error(),
			})
			fp.ResponseChannel <- res
			return
		}
		fp.ResponseChannel <- resByte
		return

	}else {
		for i := 0; i < backendsLen - 1; i++ {
			resByte, rep, _ := SendRequest(backends[i].URL(),r, reqBody)
			resultResp = append(resultResp, rep)
			result = append(result, string(resByte))
		}
	}
	if result == nil {
		//
		err := errors.New("failed to get response")
		res, _ := json.Marshal(types.ErrorResponse{
			Err:  err.Error(),
		})
		fp.ResponseChannel <- res
		return
	}
	for i := range resultResp {
		if resultResp[i].StatusCode == types.ValidCode {
			resultByte = []byte(result[i])
			break
		}
	}
	fp.ResponseChannel <- resultByte
	return

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