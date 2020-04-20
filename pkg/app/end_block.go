package app

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (c *Chain) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return c.mm.EndBlocker(ctx, req)
}
