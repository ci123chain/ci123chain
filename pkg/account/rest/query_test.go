package rest

import (
	"context"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/http"
	"testing"
)

func TestQueryHistory(t *testing.T) {
	rpc, err := http.New("tcp://localhost:36657", "/websocket")
	if err != nil {
		panic(err)
	}
	height := int64(200)
	addr := sdk.HexToAddress("0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c")
	key := types.AddressStoreKey(addr)
	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/%s/key", "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c", "accounts"),
		Height: height,
		Data:   key,
		Prove:  true,
	}
	opts := rpcclient.ABCIQueryOptions{
		Height: height,
		Prove:  req.Prove,
	}
	ctx := context.Background()
	res, err := rpc.ABCIQueryWithOptions(ctx, req.Path, req.Data, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res.Response.Value))
	fmt.Println(res.Response.ProofOps)
}