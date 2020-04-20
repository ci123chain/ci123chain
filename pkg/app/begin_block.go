package app

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func (c *Chain) BeginBlocker(ctx types.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock{

	return c.mm.BeginBlocker(ctx, req)
}

func (c *Chain) Committer(ctx types.Context) abci.ResponseCommit{

	return c.mm.Committer(ctx)
}