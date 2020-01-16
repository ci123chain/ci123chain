package server

import (
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/gateway/lbpolicy"
	"github.com/tanhuiya/ci123chain/pkg/gateway/types"

	"net/http"
)

type LBProxy struct {
	Policy types.LBPolicy
	ProxyType types.ProxyType
}


func NewLBProxy(pt types.ProxyType) *LBProxy {
	policy := lbpolicy.NewRoundPolicy()
	lbp := &LBProxy{
		ProxyType: pt,
		Policy:policy,
	}
	return lbp
}

func (lbp *LBProxy) Handle(r *http.Request, backends []types.Instance, reqBody []byte) ([]byte, error) {

	url := lbp.Policy.NextPeer(backends).URL()
	b, _, err := SendRequest(url, r, reqBody)
	if err != nil {
		return nil, errors.New("failed to get response")
	}
	return b, nil
}