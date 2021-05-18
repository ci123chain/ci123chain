package rpc

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// REST handler to get the latest block
func LatestBlockRequestHandlerFn(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getBlock(ctx, nil)
		if err != nil {
			rest.WriteErrorRes(w, types.ErrNode(types.DefaultCodespace, err).Error())
			return
		}

		rest.PostProcessResponseBare(w, ctx, output)
	}
}

func getBlock(clientCtx context.Context, height *int64) ([]byte, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(clientCtx.Context(), height)
	if err != nil {
		return nil, err
	}

	return types2.GetCodec().MarshalJSON(res)
}

// get the current blockchain height
func GetChainHeight(clientCtx context.Context) (int64, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return -1, err
	}

	status, err := node.Status(clientCtx.Context())
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// REST handler to get a block
func BlockRequestHandlerFn(clientCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorRes(w,
				"couldn't parse block height. Assumed format is '/block/{height}'.")
			return
		}

		chainHeight, err := GetChainHeight(clientCtx)
		if err != nil {
			rest.WriteErrorRes(w, "failed to parse chain height")
			return
		}

		if height > chainHeight {
			rest.WriteErrorRes(w, "requested block height is bigger then the chain length")
			return
		}

		output, err := getBlock(clientCtx, &height)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}