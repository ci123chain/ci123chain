package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	s "github.com/ci123chain/ci123chain/pkg/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"
)

func NewQuerier(k StakingKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error)  {
		switch path[0] {
		case s.QueryDelegation:
			return queryDelegation(ctx, req, k)
		case s.QueryValidatorDelegations:
			return queryAllDelegation(ctx, req, k)
		case s.QueryValidators:
			return queryValidators(ctx, req, k)
		case s.QueryValidator:
			return queryValidator(ctx, req, k)
		case s.QueryDelegatorValidators:
			return queryDelegatorValidators(ctx, req, k)
		case s.QueryDelegatorValidator:
			return queryDelegatorValidator(ctx, req, k)
		case s.QueryRedelegations:
			return queryRedelegations(ctx, req, k)
		case s.QueryDelegatorDelegations:
			return queryDelegatorDelegations(ctx, req, k)
		case s.QueryOperatorAddressSet:
			return queryOperatorAddressByConsAddresses(ctx, req, k)
		case s.QueryParameters:
			return queryParameters(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown request endpoint")
		}
	}
}

func queryDelegation(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {
	//
	var params types.QueryBondsParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	delegation, found := k.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "no delegation found")
	}

	delegationResp, err := delegationToDelegationResponse(ctx, k, delegation)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "delegation to delegationResponse failed")
	}

	res := types.StakingCodec.MustMarshalJSON(delegationResp)

	return res, nil

}

func queryAllDelegation(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {
	//
	var params types.QueryValidatorParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	delegations := k.GetValidatorDelegations(ctx, params.ValidatorAddr)
	delegationResps, err := delegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("delegation to delegationResponse failed: %v", err.Error()))
	}

	if delegationResps == nil {
		delegationResps = types.DelegationResponses{}
	}

	res := types.StakingCodec.MustMarshalJSON(delegationResps)

	return res, nil
}


func queryValidators(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {
	var params types.QueryValidatorsParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil,sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	validators := k.GetAllValidators(ctx)
	filteredVals := make([]types.Validator, 0, len(validators))

	if params.Status == sdk.BondStatusAll {
		for _, val := range validators {
			filteredVals = append(filteredVals, val)
		}
	}else {
		for _, val := range validators {
			if strings.EqualFold(val.GetStatus().String(), params.Status) {
				filteredVals = append(filteredVals, val)
			}
		}
	}

	start, end := sdk.Paginate(len(filteredVals), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		filteredVals = []types.Validator{}
	} else {
		filteredVals = filteredVals[start:end]
	}
	res := types.StakingCodec.MustMarshalJSON(filteredVals)

	return res, nil
}

func queryValidator(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {
	var params types.QueryValidatorParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	validator, found := k.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "no validator found")
	}
	res := types.StakingCodec.MustMarshalJSON(validator)

	return res, nil
}

func queryDelegatorValidators(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {
	var params types.QueryDelegatorParams

	stakingParams := k.GetParams(ctx)

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	validators, err := k.GetDelegatorValidators(ctx, params.DelegatorAddr, stakingParams.MaxValidators)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("no validator found: %v", err.Error()))
	}
	if validators == nil {
		validators = types.Validators{}
	}
	res := types.StakingCodec.MustMarshalJSON(validators)


	return res, nil
}

func queryDelegatorValidator(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {
	var params types.QueryBondsParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	validator, err := k.GetDelegatorValidator(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("get validator failed: %v", err.Error()))
	}

	res := types.StakingCodec.MustMarshalJSON(validator)

	return res, nil
}

func queryRedelegations(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {
	var params types.QueryRedelegationParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	var redels []types.Redelegation

	switch {
	case !params.DelegatorAddr.Empty() && !params.SrcValidatorAddr.Empty() && !params.DstValidatorAddr.Empty():
		redel, found := k.GetRedelegation(ctx, params.DelegatorAddr, params.SrcValidatorAddr, params.DstValidatorAddr)
		if !found {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "no redelegation")
		}

		redels = []types.Redelegation{redel}
	case params.DelegatorAddr.Empty() && !params.SrcValidatorAddr.Empty() && params.DstValidatorAddr.Empty():
		redels = k.GetRedelegationsFromSrcValidator(ctx, params.SrcValidatorAddr)
	default:
		redels = k.GetAllRedelegations(ctx, params.DelegatorAddr, params.SrcValidatorAddr, params.DstValidatorAddr)
	}

	redelResponses, err := redelegationsToRedelegationResponses(ctx, k, redels)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "redelegaion to redelegaionResponse failed")
	}

	if redelResponses == nil {
		redelResponses = types.RedelegationResponses{}
	}

	res := types.StakingCodec.MustMarshalJSON(redelResponses)

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

func delegationToDelegationResponse(ctx sdk.Context, k StakingKeeper, del types.Delegation) (types.DelegationResponse, error) {

	val, found := k.GetValidator(ctx, del.ValidatorAddress)
	if !found {
		return types.DelegationResponse{}, sdkerrors.Wrap(sdkerrors.ErrInternal, "no validator found")
	}

	return types.NewDelegationResp(
		del.DelegatorAddress,
		del.ValidatorAddress,
		del.Shares,
		sdk.NewChainCoin(val.TokensFromShares(del.Shares).TruncateInt()),
	), nil
}

func redelegationsToRedelegationResponses(
	ctx sdk.Context, k StakingKeeper, redels types.Redelegations,
) (types.RedelegationResponses, error) {

	resp := make(types.RedelegationResponses, len(redels))
	for i, redel := range redels {
		val, found := k.GetValidator(ctx, redel.ValidatorDstAddress)
		if !found {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "no validator found")
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

func queryDelegatorDelegations(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {

	var params types.QueryDelegatorParams

	err := types.StakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)
	delegationResps, err := delegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "failed to get delegation responses")
	}

	if delegationResps == nil {
		delegationResps = types.DelegationResponses{}
	}

	res := types.StakingCodec.MustMarshalJSON(delegationResps)

	return res, nil
}

func queryOperatorAddressByConsAddresses(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {

	var addressSet types.QueryOperatorAddressesParams
	err := types.StakingCodec.UnmarshalJSON(req.Data, &addressSet)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	var responseSet = make([]types.ValidatorOperatorAddressResponse, 0)
	for _, v := range addressSet.ConsAddresses {
		validator, found := k.GetValidatorByConsAddr(ctx, sdk.ToAccAddress(v))
		if !found {
			addr := types.NewValidatorOperatorAddressResponse(strings.ToUpper(hex.EncodeToString(v)), validator.OperatorAddress.String(), false)
			responseSet = append(responseSet, addr)
		}else {
			addr := types.NewValidatorOperatorAddressResponse(strings.ToUpper(hex.EncodeToString(v)), validator.OperatorAddress.String(), true)
			responseSet = append(responseSet, addr)
		}
	}
	result  := types.StakingCodec.MustMarshalJSON(responseSet)
	return result, nil
}


func queryParameters(ctx sdk.Context, req abci.RequestQuery, k StakingKeeper) ([]byte, error) {
	params := k.GetParams(ctx)
	res := types.StakingCodec.MustMarshalJSON(params)
	return res, nil
}
