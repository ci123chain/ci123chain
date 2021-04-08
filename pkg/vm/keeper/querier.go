package keeper

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	ethcmn "github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)


func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error)  {
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
		case evmtypes.QueryBloom:
			return queryBlockBloom(ctx, path, k)
		case evmtypes.QueryCode:
			return queryCode(ctx, path, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}


func queryContractInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {

	var params types.ContractInfoParams

	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	contractInfo := k.GetContractInfo(ctx, params.ContractAddress)
	if contractInfo == nil {
		return nil, nil
	}

	res := types.WasmCodec.MustMarshalBinaryBare(contractInfo)

	return res, nil
}

func queryCodeInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {

	var params types.CodeInfoParams
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}

	codeInfo := k.GetCodeInfo(ctx, params.Hash)
	if codeInfo == nil {
		return nil, nil
	}
	res := types.WasmCodec.MustMarshalBinaryBare(codeInfo)

	return res, nil
}

func queryContractState(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.ContractStateParam
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}
	//query
	args, err := types.CallData2Input(params.QueryMessage)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("invalid params: %v", err.Error()))
	}
	contractState, err := k.Query(ctx, params.ContractAddress, params.InvokerAddress, args)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("query failed: %v", err.Error()))
	}
	res := types.WasmCodec.MustMarshalJSON(contractState)

	return res, nil
}

func queryAccountContractList(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.ContractListParams
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}
	account := k.AccountKeeper.GetAccount(ctx, params.AccountAddress)
	if account == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid account")
	}
	var contractList []string
	ccstore := ctx.KVStore(k.storeKey)
	contractListBytes := ccstore.Get(types.GetAccountContractListKey(account.GetAddress()))
	if contractListBytes == nil {
		return nil, nil
	}
	err = json.Unmarshal(contractListBytes, &contractList)
	if err != nil{
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	list := types.NewContractListResponse(contractList)
	res := types.WasmCodec.MustMarshalBinaryBare(list)
	return res, nil
}

func queryContractExist(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.ContractExistParams
	err := types.WasmCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
	}
	store := ctx.KVStore(k.storeKey)
	byCode := store.Get(params.WasmCodeHash)
	if byCode != nil {
		return []byte("The contract already exists"), nil
	}
	return nil, nil
}

func queryBlockBloom(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	num, err := strconv.ParseInt(path[1], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("could not unmarshal block height: %w", err))
	}

	bloom, found := k.GetBlockBloom(ctx.WithBlockHeight(num), num)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("block bloom not found for height %d", num))
	}

	res := evmtypes.QueryBloomFilter{Bloom: bloom}
	bz, err := codec.MarshalJSONIndent(k.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal failed: %v", err.Error()))
	}

	return bz, nil
}

func queryCode(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := ethcmn.HexToAddress(path[1])
	code := keeper.GetCode(ctx, addr)
	res := evmtypes.QueryResCode{Code: code}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal failed: %v", err.Error()))
	}

	return bz, nil
}