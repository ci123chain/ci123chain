package core

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case *clienttypes.MsgCreateClient:
			res, err := k.CreateClient()
			return sdk.WrapServiceResult(ctx, res, err)
		}
	}
}