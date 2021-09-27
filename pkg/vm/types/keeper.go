package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type Keeper interface {
	EvmTxExec(ctx sdk.Context,msg sdk.Msg) (VMResult, error)
}

type VMResult interface {
	VMResult() *sdk.Result
}

