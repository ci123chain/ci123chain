package utils

import (
	"github.com/ci123chain/ci123chain/pkg/client/context"
	ibcclient "github.com/ci123chain/ci123chain/pkg/ibc/core/client"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

func QueryParms(clientCtx context.Context) (*types.Params, error) {
	path := "/custom/" + types.ModuleName + "/" + types.QueryParameters

	value, _, err := ibcclient.QueryABCI(clientCtx, path, nil, false)
	if err != nil {
		return nil, err
	}
	var p types.Params
	//err = json.Unmarshal(value, &p)
	err = types.StakingCodec.UnmarshalJSON(value, &p)
	if err != nil {
		return nil, err
	}
	return &p, err
}