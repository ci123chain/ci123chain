package context

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/types"

	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/tendermint/tendermint/crypto/merkle"

	"github.com/tendermint/tendermint/libs/bytes"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func (ctx Context) Query(path string, key bytes.HexBytes, isProve bool) ([]byte, int64, *merkle.Proof, sdk.Error) {
	var res []byte
	var height int64
	node, err := ctx.GetNode()
	if err != nil {
		return res, height, nil, transfer.ErrQueryTx(types.DefaultCodespace, err.Error())

	}

	opt := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
	}
	if isProve {
		opt.Prove = true
	}
   	result, err := node.ABCIQueryWithOptions(ctx.Context(), path, key, opt)
	if err != nil {
		return res, height, nil, transfer.ErrQueryTx(types.DefaultCodespace, err.Error())

	}

	resp := result.Response
	if !resp.IsOK() {
		return res, resp.Height, nil, transfer.ErrQueryTx(types.DefaultCodespace, resp.Log)

	}

	// verify proof

	return resp.Value, resp.Height, nil, nil
}
