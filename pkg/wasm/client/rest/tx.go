package rest

import (
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/util"
	"github.com/tanhuiya/ci123chain/pkg/wasm/types"
	"net/http"
)

func registerTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/wasm/code/install", storeCodeHandler(cliCtx)).Methods("POST")
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