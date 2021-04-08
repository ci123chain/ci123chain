package context

import (
	"fmt"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"

	"github.com/tendermint/tendermint/crypto/merkle"

	"github.com/tendermint/tendermint/libs/bytes"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func (ctx Context) Query(path string, key bytes.HexBytes, isProve bool) ([]byte, int64, *merkle.Proof, error) {
	var res []byte
	var height int64
	node, err := ctx.GetNode()
	if err != nil {
		return res, height, nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("get node failed:%v", err.Error()))

	}

	opt := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
	}
	if isProve {
		opt.Prove = true
	}
   	result, err := node.ABCIQueryWithOptions(ctx.Context(), path, key, opt)
	if err != nil {
		return res, height, nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("query failed: %v", err.Error()))

	}

	resp := result.Response
	if !resp.IsOK() {
		return res, resp.Height, nil, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("query failed, got: %v", resp))

	}

	// verify proof

	return resp.Value, resp.Height, nil, nil
}
