package server

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/gateway/lbpolicy"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"net/http"
)

type LBProxy struct {
	Policy types.LBPolicy
	ProxyType types.ProxyType
	ResponseChannel chan []byte
}


func NewLBProxy(pt types.ProxyType) *LBProxy {
	policy := lbpolicy.NewRoundPolicy()
	lbp := &LBProxy{
		ProxyType: pt,
		Policy:policy,
		ResponseChannel:make(chan []byte),
	}
	return lbp
}

func (lbp *LBProxy) Handle(r *http.Request, backends []types.Instance, RequestParams map[string]string) []byte {

	url := lbp.Policy.NextPeer(backends).URL()
	b, _, err := SendRequest(url, r, RequestParams)
	if err != nil {
		res, _ := json.Marshal(types.ErrorResponse{
			Ret: 0,
			Message:  err.Error(),
		})
		lbp.ResponseChannel <- res
		return res
	}
	lbp.ResponseChannel <- b
	return b
}

func (lbp *LBProxy) Response() *chan []byte {
	return &lbp.ResponseChannel
}