package server

import (
	//"encoding/json"
	//"fmt"

	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"net/http"
)


type ConcretProxy struct {
	ProxyType types.ProxyType
	//Backends []types.Instance
}


func NewConcretProxy(pt types.ProxyType) ConcretProxy {

	cp := ConcretProxy{
		ProxyType: pt,
	}
	//cp.ConfigServerPool(list)

	return cp
}


func (cp ConcretProxy) Handle(w http.ResponseWriter, r *http.Request, backends []types.Instance) {
	//

	//backendsLen := len(backends)
	//var resultResp []types.ResultRep
	/*cli := &http.Client{}
	body := make([]byte, 0)
	reqUrl := "http://" + remote_addr + r.URL.Path

	req2, err := http.NewRequest(r.Method, reqUrl, strings.NewReader(string(body)))
	if err != nil {
		io.WriteString(w, "Request Error")
		return
	}
	// set request content type
	contentType := r.Header.Get("Content-Type")
	req2.Header.Set("Content-Type", contentType)
	// request
	rep2, err := cli.Do(req2)
	if err != nil {
		io.WriteString(w, "Not Found!")
		return
	}
	b, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		panic(err)
	}*/

	/*
	if resLen == 1 {
		var result types.ResultRep
		byte, err := SendRequest(cp.Backends[0].URL(), r)
		if err != nil {
			//
		}
		err = json.Unmarshal(byte, result)
		resultResp = append(resultResp, result)
	}else {
		for i := 0; i < resLen - 1; i++ {
			//
			var result types.ResultRep
			byte, err := SendRequest(cp.Backends[i].URL(),r)
			if err != nil {
				//
			}
			err = json.Unmarshal(byte, &result)
			if err != nil {
				//
			}
			resultResp = append(resultResp, result)
		}
	}
	resultByte, err := json.Marshal(resultResp)
	if err != nil {
		//
	}
	w.Write(resultByte)
	*/

	//------
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
	allResult, err := json.Marshal(ResultResponse)
	if err != nil {
		//
	}
	_, err = w.Write(allResult)
	if err != nil {
		//
	}
	clasterTaskPool.Stop()
	*/

	b, rep, err := SendRequest(backends[0].URL(), r)
	if err != nil {
		panic(err)
	}
	/*
	var res ciRes
	err = json.Unmarshal(b, &res)
	if err != nil {
		fmt.Println("222")
		fmt.Println(string(b))
	}
	fmt.Println(res.Data.Balance)

	//fmt.Println(res)
	w.Write(b)
	*/

	//
	var resultCi ciResult
	var resCi ciRes
	err = json.Unmarshal(b, &resultCi)
	if err != nil {
		err2 := json.Unmarshal(b, &resCi)
		if err2 != nil {
			w.Write([]byte("failed to unmarshal"))
		}
		var data = BalanceData{resCi.Data.Balance}
		data_byte,err2 := json.Marshal(data)
		if err != nil {
			//
		}
		var resData = ciResult{
			Ret:     0,
			Data:    string(data_byte),
			Message: "",
		}
		response_byte, err2 := json.Marshal(resData)
		if err2 != nil {
			//
		}
		copyHeader(w.Header(), rep.Header)
		w.WriteHeader(rep.StatusCode)
		w.Write(response_byte)
	}
	w.Write(b)
}
/*
func (cp ConcretProxy) AddBackend(backend types.Instance) {
	cp.Backends = append(cp.Backends, backend)
}

func (cp ConcretProxy)ConfigServerPool(tokens []string)  {
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		cp.AddBackend(backend.NewBackEnd(serverUrl, true, nil))
	}
}
*/

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}