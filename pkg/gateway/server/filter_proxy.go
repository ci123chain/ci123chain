package server

import (
	"encoding/json"
	"fmt"

	"github.com/tanhuiya/ci123chain/pkg/gateway/types"
	"net/http"
)

type FilterProxy struct {
	ProxyType types.ProxyType
	//Backends []types.Instance
}

func NewFilterProxy(pt types.ProxyType) FilterProxy {

	fp :=  FilterProxy{
		ProxyType:pt,
	}
	//fp.ConfigServerPool(list)
	return fp
}

func (fp FilterProxy) Handle(w http.ResponseWriter,r *http.Request, backends []types.Instance) {
	//
	backendsLen := len(backends)
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
		byte, err := SendRequest(fp.Backends[0].URL(), r)
		if err != nil {
			//
		}
		w.Write(byte)

	}else {
		for i := 0; i < resLen - 1; i++ {
			//
			var result types.ResultRep
			byte, err := SendRequest(fp.Backends[i].URL(),r)
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
	for i, _ := range resultResp {
		var byte []byte
		var err error
		if resultResp[i].Code == 200 {
			byte, err = json.Marshal(resultResp[i])
			if err != nil {
				//
			}
			break
		}
		w.Write(byte)
	}
	*/
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
	w.Write(byte)

	clasterTaskPool.Stop()
}
/*
func (fp FilterProxy) AddBackend(backend types.Instance) {
	fp.Backends = append(fp.Backends, backend)
}

func (fp FilterProxy)ConfigServerPool(tokens []string)  {
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		fp.AddBackend(backend.NewBackEnd(serverUrl, true, nil))
	}
}
*/

func filter(resultResponse []types.ResultRep, b bool) (byte []byte){
	for i := range resultResponse {
		if resultResponse[i].Code != types.ValidCode {
			b = false
		}else {
			b = true
			byte, _ = json.Marshal(resultResponse[i])
			return
		}
	}
	byte, _ = json.Marshal(resultResponse[0])
	return
}