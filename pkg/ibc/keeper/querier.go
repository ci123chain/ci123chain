package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	QueryState = "state"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper IBCKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryState:
			return queryResolve(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

// nolint: unparam
func queryResolve(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper IBCKeeper) ([]byte, sdk.Error) {
	if path[0] != types.StateReady {
		return nil, transfer.ErrCheckParams(types.DefaultCodespace, "Parameter State Error")
	}

	value := keeper.GetFirstReadyIBCMsg(ctx)
	if value == nil {
		return nil, transfer.ErrQueryTx(types.DefaultCodespace, "No ready ibc tx found")
	}

	retbz, err := types.IbcCdc.MarshalBinaryLengthPrefixed(*value)
	if err != nil {
		return nil, types.ErrFailedMarshal(types.DefaultCodespace, err.Error())
	}
	return retbz, nil
}
/*
func queryAccountNonce(ctx sdk.Context, path []string, req abci.RequestQuery, keeper IBCKeeper) ([]byte, sdk.Error) {
	accountAddr := path[0]
	address := sdk.HexToAddress(accountAddr)
	account := keeper.AccountKeeper.GetAccount(ctx, address)
	if account == nil {
		return nil, transfer.ErrQueryTx(types.DefaultCodespace, "account is not exist")
	}
	nonce := account.GetSequence()
	retbz, err := types.IbcCdc.MarshalBinaryLengthPrefixed(nonce)
	if err != nil {
		return nil, types.ErrFailedMarshal(types.DefaultCodespace, err.Error())
	}
	return retbz, nil
}*/