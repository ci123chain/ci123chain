package utils

import (
	"github.com/ci123chain/ci123chain/pkg/client/context"
	ibcclient "github.com/ci123chain/ci123chain/pkg/ibc/core/client"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

func QueryParms(clientCtx context.Context) (p *types.Params, err error) {
	path := "/custom/" + types.ModuleName + "/" + types.QueryParameters

	value, _, err := ibcclient.QueryABCI(clientCtx, path, nil, false)
	if err != nil {
		return nil, err
	}
	err = types.StakingCodec.UnmarshalJSON(value, p)
	if err != nil {
		return nil, err
	}
	return p, err
}