package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	QueryState = "state"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper IBCKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryState:
			return queryResolve(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}

// nolint: unparam
func queryResolve(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper IBCKeeper) ([]byte, error) {
	if path[0] != types.StateReady {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, "Paramter state error")
	}

	value := keeper.GetFirstReadyIBCMsg(ctx)
	if value == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, "No ready ibc tx found")
	}

	retbz, err := types.IbcCdc.MarshalBinaryLengthPrefixed(*value)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal failed: %v", err.Error()))
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