package app

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
)


func (c *Chain) BeginBlocker(ctx types.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock{

	return c.mm.BeginBlocker(ctx, req)
}
