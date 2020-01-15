package server

import (
	"github.com/tanhuiya/ci123chain/pkg/gateway/lbpolicy"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"

	"net/http"
)

type LBProxy struct {
	Policy types.LBPolicy
	ProxyType types.ProxyType
	//Backends []types.Instance
}


func NewLBProxy(pt types.ProxyType) LBProxy {
	policy := lbpolicy.NewRoundPolicy()
	lbp := LBProxy{
		ProxyType: pt,
		Policy:policy,
	}
	//lbp.ConfigServerPool(list)
	return lbp
}

func (lbp LBProxy) Handle(w http.ResponseWriter, r *http.Request, backends []types.Instance) {

	/*attempts := gw.GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}*/
	//

	rw := w
	peer := lbp.Policy.NextPeer(backends)
	peer.Proxy().ServeHTTP(rw ,r)
/*
	b, err := SendRequest(backends[0].URL(), r)
	if err != nil {
		panic(err)
	}
	//1.如果第一个unmarshal成功
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
		w.Write(response_byte)
	}
	*/
}
/*
func (lbp LBProxy) AddBackend(backend types.Instance) {
	lbp.Backends = append(lbp.Backends, backend)
}

func (lbp LBProxy)ConfigServerPool(tokens []string)  {
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		/*proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := GetRetryFromContext(request)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			// after 3 retries, mark this backend as down
			serverPool.MarkBackendStatus(serverUrl, false)

			// if the same request routing for few attempts with different backends, increase the count
			attempts := GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			AllHandle(writer, request.WithContext(ctx))
		}*/
/*
		lbp.AddBackend(backend.NewBackEnd(serverUrl, true, proxy))
	}
}
*/