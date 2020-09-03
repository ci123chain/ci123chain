package common

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)


func QueryValidatorCommission(cliCtx context.Context,queryRouter string, val sdk.AccAddress) ([]byte, error) {

<<<<<<< HEAD
	res, _, _, err := cliCtx.Query(fmt.Sprintf("/custom/%s/%s", queryRouter, types.QueryValidatorCommission),
		cliCtx.Cdc.MustMarshalJSON(types.NewQueryValidatorCommissionParams(val)), false)
=======
	res, _, err := cliCtx.Query(fmt.Sprintf("/custom/%s/%s", queryRouter, types.QueryValidatorCommission),
		cliCtx.Cdc.MustMarshalJSON(types.NewQueryValidatorCommissionParams(val)))
>>>>>>> mint

	return res, err
}

func QueryDelegationRewards(cliCtx context.Context, queryRouter string, val,  del sdk.AccAddress) ([]byte, error) {
<<<<<<< HEAD
	res, _, _, err := cliCtx.Query(fmt.Sprintf("/custom/%s/%s", queryRouter, types.QueryDelegationRewards),
		cliCtx.Cdc.MustMarshalJSON(types.NewQueryDelegationRewardsParams(del, val)), false)
=======
	res, _, err := cliCtx.Query(fmt.Sprintf("/custom/%s/%s", queryRouter, types.QueryDelegationRewards),
		cliCtx.Cdc.MustMarshalJSON(types.NewQueryDelegationRewardsParams(del, val)))
>>>>>>> mint

	return res, err
}