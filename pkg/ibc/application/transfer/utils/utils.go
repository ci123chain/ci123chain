package utils

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/ibc"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	ibcclient "github.com/ci123chain/ci123chain/pkg/ibc/core/client"
)

func QueryDenomTraces(clientCtx context.Context, offset, limit uint64) (*types.QueryDenomTracesResponse, error) {
	path := "/custom/" + ibc.ModuleName + "/" + types.QueryDenomTraces
	req := &types.QueryDenomTracesRequest{
		Pagination: &pagination.PageRequest{
			Offset:     offset,
			Limit:      limit,
			CountTotal: true,
		},
	}
	key := clientCtx.Cdc.MustMarshalJSON(req)
	value, _, err := ibcclient.QueryABCI(clientCtx, path, key, false)
	if err != nil {
		return nil, err
	}

	var resp types.QueryDenomTracesResponse
	clientCtx.Cdc.MustUnmarshalJSON(value, &resp)
	return &resp, nil
}