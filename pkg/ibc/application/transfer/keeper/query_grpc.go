package keeper

import (
	store2 "github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DenomTraces implements the Query/DenomTraces gRPC method
func (q Keeper) DenomTraces(ctx sdk.Context, req *types.QueryDenomTracesRequest) (*types.QueryDenomTracesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	//ctx := sdk.UnwrapSDKContext(c)

	traces := types.Traces{}
	store := store2.NewPrefixStore(ctx.KVStore(q.storeKey), types.DenomTraceKey)

	pageRes, err := pagination.Paginate(store, req.Pagination, func(_, value []byte) error {
		result, err := q.UnmarshalDenomTrace(value)
		if err != nil {
			return err
		}

		traces = append(traces, result)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryDenomTracesResponse{
		DenomTraces: traces.Sort(),
		Pagination:  pageRes,
	}, nil
}
