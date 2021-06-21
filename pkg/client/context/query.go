package context

import (
	"context"
	"fmt"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"

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


// QueryABCI performs a query to a Tendermint node with the provide RequestQuery.
// It returns the ResultQuery obtained from the query.
func (ctx Context) QueryABCI(req abci.RequestQuery) (abci.ResponseQuery, error) {
	return ctx.queryABCI(req)
}


func (ctx Context) queryABCI(req abci.RequestQuery) (abci.ResponseQuery, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return abci.ResponseQuery{}, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
		Prove:  req.Prove,
	}

	result, err := node.ABCIQueryWithOptions(context.Background(), req.Path, req.Data, opts)
	if err != nil {
		return abci.ResponseQuery{}, err
	}

	if !result.Response.IsOK() {
		return abci.ResponseQuery{}, errors.New(result.Response.Log)
	}

	// data from trusted node or subspace query doesn't need verification
	if !opts.Prove || !isQueryStoreWithProof(req.Path) {
		return result.Response, nil
	}

	return result.Response, nil
}


// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType must be "store" and subpath must be "key" to require a proof.
func isQueryStoreWithProof(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}

	paths := strings.SplitN(path[1:], "/", 3)

	switch {
	case len(paths) != 3:
		return false
	case paths[0] != "store":
		return false
	case RequireProof("/" + paths[2]):
		return true
	}

	return false
}

// RequireProof returns whether proof is required for the subpath.
func RequireProof(subpath string) bool {
	// XXX: create a better convention.
	// Currently, only when query subpath is "/key", will proof be included in
	// response. If there are some changes about proof building in iavlstore.go,
	// we must change code here to keep consistency with iavlStore#Query.
	return subpath == "/key"
}
