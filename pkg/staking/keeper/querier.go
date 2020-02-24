package keeper

import (
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
	s "github.com/tanhuiya/ci123chain/pkg/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"
)

func NewQuerier(k StakingKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error)  {
		switch path[0] {
		case s.QueryDelegation://
			return queryDelegation(ctx, req, k)
		case s.QueryValidatorDelegations://
			return queryAllDelegation(ctx, req, k)
		case s.QueryValidators://
			return queryValidators(ctx, req, k)
		case s.QueryValidator://
			return queryValidator(ctx, req, k)
		case s.QueryDelegatorValidators://
			return queryDelegatorValidators(ctx, req, k)
		case s.QueryDelegatorValidator://
			return queryDelegatorValidator(ctx, req, k)
		case s.QueryRedelegations:
			return queryRedelegations(ctx, req, k)
		case s.QueryDelegatorDelegations:
			return queryDelegatorDelegations(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryDelegation(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, sdk.Error) {
	//
	var params types.QueryBondsParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	delegation, found := k.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return nil, sdk.ErrNoDelegation("no delegation")
	}

	delegationResp, err := delegationToDelegationResponse(ctx, k, delegation)
	if err != nil {
		return nil, sdk.ErrInternal("delegation to delegationResponse failed")
	}

	res := types.StakingCodec.MustMarshalJSON(delegationResp)
	/*if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}*/

	return res, nil

	//return nil, nil
}

func queryAllDelegation(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, sdk.Error) {
	//
	var params types.QueryValidatorParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	delegations := k.GetValidatorDelegations(ctx, params.ValidatorAddr)
	delegationResps, err := delegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return nil, sdk.ErrInternal("failed get delegation responses")
	}

	if delegationResps == nil {
		delegationResps = types.DelegationResponses{}
	}

	res := types.StakingCodec.MustMarshalJSON(delegationResps)

	return res, nil
}


func queryValidators(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorsParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	validators := k.GetAllValidators(ctx)
	fmt.Println(validators)
	filteredVals := make([]types.Validator, 0, len(validators))

	for _, val := range validators {
		if strings.EqualFold(val.GetStatus().String(), params.Status) {
			filteredVals = append(filteredVals, val)
		}
	}

	start, end := sdk.Paginate(len(filteredVals), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		filteredVals = []types.Validator{}
	} else {
		filteredVals = filteredVals[start:end]
	}
	res:= types.StakingCodec.MustMarshalJSON(filteredVals)

	return res, nil
}

func queryValidator(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	validator, found := k.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return nil, sdk.ErrNoValidatorFound("no validator found")
	}
	res := types.StakingCodec.MustMarshalJSON(validator)

	/*res, err := app.MarshalJSONIndent(types.StakingCodec, validator)
	if err != nil {
		return nil, sdk.ErrInternal("marshal failed")
	}*/

	return res, nil
}

func queryDelegatorValidators(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, sdk.Error) {
	var params types.QueryDelegatorParams

	stakingParams := k.GetParams(ctx)

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	validators, err := k.GetDelegatorValidators(ctx, params.DelegatorAddr, stakingParams.MaxValidators)
	if err != nil {
		return nil, sdk.ErrNoValidatorFound("no validator found")
	}
	if validators == nil {
		validators = types.Validators{}
	}
	res := types.StakingCodec.MustMarshalJSON(validators)

	/*res, err := app.MarshalJSONIndent(types.StakingCodec, validators)
	if err != nil {
		return nil, sdk.ErrInternal("marshal failed")
	}*/

	return res, nil
}

func queryDelegatorValidator(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, sdk.Error) {
	var params types.QueryBondsParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	validator, err := k.GetDelegatorValidator(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if err != nil {
		return nil, sdk.ErrInternal("get validator failed")
	}

	res := types.StakingCodec.MustMarshalJSON(validator)
	/*res, err := app.MarshalJSONIndent(types.StakingCodec, validator)
	if err != nil {
		return nil, sdk.ErrInternal("marshal failed")
	}*/

	return res, nil
}

func queryRedelegations(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, sdk.Error) {
	var params types.QueryRedelegationParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	var redels []types.Redelegation

	switch {
	case !params.DelegatorAddr.Empty() && !params.SrcValidatorAddr.Empty() && !params.DstValidatorAddr.Empty():
		redel, found := k.GetRedelegation(ctx, params.DelegatorAddr, params.SrcValidatorAddr, params.DstValidatorAddr)
		if !found {
			return nil, sdk.ErrNoRedelegation("no redelegation")
		}

		redels = []types.Redelegation{redel}
	case params.DelegatorAddr.Empty() && !params.SrcValidatorAddr.Empty() && params.DstValidatorAddr.Empty():
		redels = k.GetRedelegationsFromSrcValidator(ctx, params.SrcValidatorAddr)
	default:
		redels = k.GetAllRedelegations(ctx, params.DelegatorAddr, params.SrcValidatorAddr, params.DstValidatorAddr)
	}

	redelResponses, err := redelegationsToRedelegationResponses(ctx, k, redels)
	if err != nil {
		return nil, sdk.ErrInternal("redelegaion to redelegaionResponse failed")
	}

	if redelResponses == nil {
		redelResponses = types.RedelegationResponses{}
	}

	res := types.StakingCodec.MustMarshalJSON(redelResponses)
	/*res, err := app.MarshalJSONIndent(types.StakingCodec, redelResponses)
	if err != nil {
		return nil, sdk.ErrInternal("marshal failed")
	}*/

	return res, nil
}


func delegationsToDelegationResponses(
	ctx sdk.Context, k StakingKeeper, delegations types.Delegations,
) (types.DelegationResponses, error) {

	resp := make(types.DelegationResponses, len(delegations))
	for i, del := range delegations {
		delResp, err := delegationToDelegationResponse(ctx, k, del)
		if err != nil {
			return nil, err
		}

		resp[i] = delResp
	}

	return resp, nil
}

func delegationToDelegationResponse(ctx sdk.Context, k StakingKeeper, del types.Delegation) (types.DelegationResponse, sdk.Error) {

	val, found := k.GetValidator(ctx, del.ValidatorAddress)
	if !found {
		return types.DelegationResponse{}, sdk.ErrNoValidatorFound("no validator found")
	}

	return types.NewDelegationResp(
		del.DelegatorAddress,
		del.ValidatorAddress,
		del.Shares,
		sdk.NewCoin(val.TokensFromShares(del.Shares).TruncateInt()),
	), nil
}

func redelegationsToRedelegationResponses(
	ctx sdk.Context, k StakingKeeper, redels types.Redelegations,
) (types.RedelegationResponses, error) {

	resp := make(types.RedelegationResponses, len(redels))
	for i, redel := range redels {
		val, found := k.GetValidator(ctx, redel.ValidatorDstAddress)
		if !found {
			return nil, sdk.ErrNoValidatorFound("no validator found")
		}

		entryResponses := make([]types.RedelegationEntryResponse, len(redel.Entries))
		for j, entry := range redel.Entries {
			entryResponses[j] = types.NewRedelegationEntryResponse(
				entry.CreationHeight,
				entry.CompletionTime,
				entry.SharesDst,
				entry.InitialBalance,
				val.TokensFromShares(entry.SharesDst).TruncateInt(),
			)
		}

		resp[i] = types.NewRedelegationResponse(
			redel.DelegatorAddress,
			redel.ValidatorSrcAddress,
			redel.ValidatorDstAddress,
			entryResponses,
		)
	}

	return resp, nil
}

func queryDelegatorDelegations(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, sdk.Error) {

	var params types.QueryDelegatorParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)
	delegationResps, err := delegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return nil, sdk.ErrInternal("failed to get delegation responses")
	}

	if delegationResps == nil {
		delegationResps = types.DelegationResponses{}
	}

	res := types.StakingCodec.MustMarshalJSON(delegationResps)

	return res, nil
}