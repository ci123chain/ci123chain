package server

import (
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"net/http"
)

type FilterProxy struct {
	ProxyType types.ProxyType
}

func NewFilterProxy(pt types.ProxyType) *FilterProxy {

	fp :=  &FilterProxy{
		ProxyType:pt,
	}
	return fp
}

func (fp *FilterProxy) Handle(r *http.Request, backends []types.Instance) ([]byte, error) {

	backendsLen := len(backends)
	var resultResp []*http.Response
	var resultByte []byte
	var result []string

	if backendsLen == 1 {
		byte, _, err := SendRequest(backends[0].URL(), r)
		if err != nil {
			//
			return nil, errors.New("failed to get response")
		}
		return byte, nil

	}else {
		for i := 0; i < backendsLen - 1; i++ {
			byte, rep2, err := SendRequest(backends[i].URL(),r)
			if err != nil {
				//
				//return nil, errors.New("failed to get response")
			}
			resultResp = append(resultResp, rep2)
			result = append(result, string(byte))
		}
	}
	for i, _ := range resultResp {
		if resultResp[i].StatusCode == types.ValidCode {
			resultByte = []byte(result[i])
			break
		}
	}
	return resultByte, nil

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