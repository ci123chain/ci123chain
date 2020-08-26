package rest

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	"net/http"
)

const CAN_MIGRATE string = `{"method":"canMigrate","args": [""]}`


func instantiateContractHandler(cliCtx context.Context,w http.ResponseWriter, r *http.Request) {
	txByte, err := buildInstantiateContractMsg(r)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"data error"))
		return
	}
	res, err := cliCtx.BroadcastTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}

func executeContractHandler(cliCtx context.Context,w http.ResponseWriter, r *http.Request) {
	contractAddr := r.FormValue("contractAddress")
	contractAddress := sdk.HexToAddress(contractAddr)
	params := types.NewQueryContractInfoParams(contractAddress)
	bz, Er := cliCtx.Cdc.MarshalJSON(params)
	if Er != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
		return
	}

	res, _, _, _:= cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractInfo, bz, false)
	if res == nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "contract does not exist or get contract error"))
		return
	}

	txByte, err := buildExecuteContractMsg(r)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"data error"))
		return
	}
	resTx, err := cliCtx.BroadcastTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, resTx)
}

func migrateContractHandler(cliCtx context.Context,w http.ResponseWriter, r *http.Request) {
	sender := cliCtx.FromAddr
	contractAddr := r.FormValue("contractAddress")
	contractAddress := sdk.HexToAddress(contractAddr)
	queryParam := []byte(CAN_MIGRATE)
	params := types.NewContractStateParam(contractAddress, sender, queryParam)
	bz, Er := cliCtx.Cdc.MarshalJSON(params)
	if Er != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
		return
	}

	res, _, _, _ := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractState, bz, false)
	var contractState types.ContractState
	cliCtx.Cdc.MustUnmarshalJSON(res, &contractState)
	if contractState.Result != "true" {
		rest.WriteErrorRes(w, sdk.ErrInternal("No permissions to migrate contracts"))
		return
	}
	txByte, err := buildMigrateContractMsg(r)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"data error"))
		return
	}
	resTx, err := cliCtx.BroadcastTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, resTx)
}