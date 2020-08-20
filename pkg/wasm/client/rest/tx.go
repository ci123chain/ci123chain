package rest

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	"github.com/gorilla/mux"
	"net/http"
)

const CAN_MIGRATE string = `{"method":"canMigrate","args": [""]}`
func registerTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/wasm/contract/init", instantiateContractHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/execute", executeContractHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/migrate", migrateContractHandler(cliCtx)).Methods("POST")
}

func instantiateContractHandler(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		async := r.FormValue("async")
		ok, err := util.CheckBool(async)  //default async
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"error async"))
			return
		}

		txByte, err := buildInstantiateContractMsg(r)
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"data error"))
			return
		}

		if ok {
			//async
			res, err := cliCtx.BroadcastTxAsync(txByte)
			if err != nil {
				rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
				return
			}
			rest.PostProcessResponseBare(w, cliCtx, res)
		}else {
			//sync
			res, err := cliCtx.BroadcastSignedData(txByte)
			if err != nil {
				rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
				return
			}
			rest.PostProcessResponseBare(w, cliCtx, res)
		}
	}
}

func executeContractHandler(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		async := r.FormValue("async")
		ok, err := util.CheckBool(async)  //default async
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"error async"))
			return
		}

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
		if ok {
			//async
			res, err := cliCtx.BroadcastTxAsync(txByte)
			if err != nil {
				rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
				return
			}
			rest.PostProcessResponseBare(w, cliCtx, res)
		}else {
			//sync
			res, err := cliCtx.BroadcastSignedData(txByte)
			if err != nil {
				rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
				return
			}
			rest.PostProcessResponseBare(w, cliCtx, res)
		}
	}
}

func migrateContractHandler(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		async := r.FormValue("async")
		ok, err := util.CheckBool(async)  //default async
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"error async"))
			return
		}
		from := r.FormValue("from")
		sender := sdk.HexToAddress(from)
		contractAddr := r.FormValue("contractAddress")
		contractAddress := sdk.HexToAddress(contractAddr)
		queryParam := []byte(CAN_MIGRATE)
		params := types.NewContractStateParam(contractAddress, sender, queryParam)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractState, bz, false)
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

		if ok {
			//async
			res, err := cliCtx.BroadcastTxAsync(txByte)
			if err != nil {
				rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
				return
			}
			rest.PostProcessResponseBare(w, cliCtx, res)
		}else {
			//sync
			res, err := cliCtx.BroadcastSignedData(txByte)
			if err != nil {
				rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
				return
			}
			rest.PostProcessResponseBare(w, cliCtx, res)
		}
	}
}