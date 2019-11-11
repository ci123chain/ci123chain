package context

import (
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func (ctx Context) Query(path string, key common.HexBytes) (res []byte, height int64, err error) {
	node, err := ctx.GetNode()
	if err != nil {
		return res, height, err
	}

	opt := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
	}
	result, err := node.ABCIQueryWithOptions(path, key, opt)
	if err != nil {
		return res, height, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, resp.Height, errors.New(resp.Log)
	}

	// todo verify proof

	return resp.Value, resp.Height, nil
}
