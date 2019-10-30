package rpc

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/abci/version"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tendermint/tendermint/p2p"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"net/http"
)

type NodeInfoResponse struct {
	p2p.DefaultNodeInfo `json:"node_info"`

	ApplicationVersion   version.Info `json:"application_version"`
}

func NodeInfoRequestHandlerFn(ctx context.Context) http.HandlerFunc  {
	return func(w http.ResponseWriter, request *http.Request) {
		status, err := getNodeStatus(ctx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		resp := NodeInfoResponse{
			DefaultNodeInfo: status.NodeInfo,
			ApplicationVersion: version.NewInfo(),
		}
		rest.PostProcessResponseBare(w, ctx, resp)
	}
}

func getNodeStatus(ctx context.Context) (*ctypes.ResultStatus, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}
	return node.Status()
}