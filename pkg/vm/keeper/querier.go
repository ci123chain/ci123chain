package keeper

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
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
		case types.QueryContractExist:
			return queryContractExist(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint")
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

	codeInfo := k.GetCodeInfo(ctx, params.Hash)
	if codeInfo == nil {
		return nil, nil
	}
	res := types.WasmCodec.MustMarshalBinaryBare(codeInfo)

	return res, nil
}

func queryContractState(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.ContractStateParam
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}
	//query
	contractState, err := k.Query(ctx, params.ContractAddress, params.InvokerAddress, params.QueryMessage)
	if err != nil {
		return nil, types.ErrQueryFailed(types.DefaultCodespace, err)
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
	var contractList []string
	contractListBytes := store.Get(types.GetAccountContractListKey(account.GetAddress()))
	if contractListBytes == nil {
		return nil, nil
	}
	err = json.Unmarshal(contractListBytes, &contractList)
	if err != nil{
		return nil, sdk.ErrInternal(err.Error())
	}
	list := types.NewContractListResponse(contractList)
	res := types.WasmCodec.MustMarshalBinaryBare(list)
	return res, nil
}

func queryContractExist(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.ContractExistParams
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}
	store := ctx.KVStore(k.storeKey)
	byCode := store.Get(params.WasmCodeHash)
	if byCode != nil {
		return []byte("The contract already exists"), nil
	}
	return nil, nil
}