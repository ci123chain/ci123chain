package rpc

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/abci/version"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/types"
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
			rest.WriteErrorRes(w, types.ErrNode(types.DefaultCodespace, err).Error())
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
	return node.Status(ctx.Context())
}