package keeper

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/wasm/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error)  {
		switch path[0] {
		case types.QueryContractInfo:
			return queryContractInfo(ctx, req, k)
		case types.QueryCodeInfo:
			return queryCodeInfo(ctx, req, k)
		case types.QueryContractState:
			return queryContractState(ctx, req, k)
		case types.QueryContractList:
			return queryAccountContractList(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}


func queryContractInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {

	var params types.ContractInfoParams

	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	contractInfo := k.GetContractInfo(ctx, params.ContractAddress)
	if contractInfo == nil {
		return nil, nil
	}
	res := types.WasmCodec.MustMarshalBinaryBare(contractInfo)

	return res, nil
}

func queryCodeInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {

	var params types.CodeInfoParams
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	codeInfo := k.GetCodeInfo(ctx, params.ID)
	if codeInfo == nil {
		return nil, nil
	}
	res := types.WasmCodec.MustMarshalBinaryBare(codeInfo)

	return res, nil
}

func queryContractState(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.ContractInfoParams
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	contractState, err := k.Query(ctx, params.ContractAddress)
	if err != nil {
		return nil, sdk.ErrInternal("get contract state failed")
	}
	res := types.WasmCodec.MustMarshalJSON(contractState)

	return res, nil
}

func queryAccountContractList(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.ContractListParams
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}
	account := k.AccountKeeper.GetAccount(ctx, params.AccountAddress)
	if account == nil {
		return nil, sdk.ErrInternal("account doesn't exists")
	}
	contractList := account.GetContractList()
	list := types.NewContractListResponse(contractList)
	res := types.WasmCodec.MustMarshalBinaryBare(list)
	return res, nil
}