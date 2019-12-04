package context

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/types"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func (ctx Context) Query(path string, key common.HexBytes) ([]byte, int64, sdk.Error) {
	var res []byte
	var height int64
	node, err := ctx.GetNode()
	if err != nil {
		return res, height, transfer.ErrQueryTx(types.DefaultCodespace, err.Error())
	}

	opt := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
	}
   	result, err := node.ABCIQueryWithOptions(path, key, opt)
	if err != nil {
		return res, height, transfer.ErrQueryTx(types.DefaultCodespace, err.Error())
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, resp.Height, transfer.ErrQueryTx(types.DefaultCodespace, resp.Log)
	}

	// todo verify proof

	return resp.Value, resp.Height, nil
}
