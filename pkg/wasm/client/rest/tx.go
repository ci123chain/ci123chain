package rest

import (
	"github.com/gorilla/mux"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ci123chain/ci123chain/pkg/wasm/keeper"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	"net/http"
)

func registerTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/wasm/contract/install", storeCodeHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/init", instantiateContractHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/execute", executeContractHandler(cliCtx)).Methods("POST")
}

func storeCodeHandler(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		async := r.FormValue("async")
		ok, err := util.CheckBool(async)  //default async
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"error async"))

			return
		}

		wasmCode, err := getWasmCode(r)
		if err != nil || wasmCode == nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "get wasmCode error"))
			return
		}
		//checkContractIsExist
		hash := keeper.MakeCodeHash(wasmCode)
		params := types.NewContractExistParams(hash)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "marshal failed"))
			return
		}
		_, _, _, err = cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractExist, bz)
		if err != nil {
			rest.WriteErrorRes(w, err.(sdk.Error))
			return
		}

		txByte, err := buildStoreCodeMsg(r)
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
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

func instantiateContractHandler(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		async := r.FormValue("async")
		ok, err := util.CheckBool(async)  //default async
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"error async"))
			return
		}

		//check codeID
		codeHash := r.FormValue("codeHash")
		params := types.NewQueryCodeInfoParams(codeHash)
		bz, err := cliCtx.Cdc.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, _, _ := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryCodeInfo, bz)
		if res == nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"codeHash does not exists"))
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

		res, _, _, _:= cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractInfo, bz)
		if res == nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "contract does not exist"))
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