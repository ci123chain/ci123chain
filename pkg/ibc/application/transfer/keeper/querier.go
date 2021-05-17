package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {

		case types.QueryDenomTraces:
			return queryDenomTraces(ctx, req, k)
		default:
			return nil, sdkerrors.ErrUnknownRequest
		}
	}
}

// ClientState implements the IBC QueryServer interface
func queryDenomTraces(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var reqDenomTraces types.QueryDenomTracesRequest

	if err := types.IBCTransferCdc.UnmarshalJSON(req.Data, &reqDenomTraces); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	resp, err := keeper.DenomTraces(ctx, &reqDenomTraces)

	respbz := types.IBCTransferCdc.MustMarshalJSON(resp)
	return respbz, err
}
