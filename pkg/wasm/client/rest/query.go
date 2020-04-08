package rest

import (
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"github.com/tanhuiya/ci123chain/pkg/wasm/types"
	"net/http"
	"strconv"
)

func registerQueryRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/wasm/codeSearch/list", listCodesHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/codeSearch", queryCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/account/contractsList", listContractsByCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contractSearch", queryContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contractSearch/state", queryContractStateAllHandlerFn(cliCtx)).Methods("POST")
}

func listCodesHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest.WriteErrorRes(w, sdk.ErrInternal("Implement me"))
	}
}


func queryCodeHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		codeId := r.FormValue("codeID")
		codeID, err := strconv.ParseUint(codeId, 10, 64)
		if err != nil {
			//rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx, ok, Err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok {
			rest.WriteErrorRes(w, Err)
			return
		}
		params := types.NewQueryCodeInfoParams(codeID)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, Err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryCodeInfo, bz)
		if Err != nil {
			rest.WriteErrorRes(w, Err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "no excepted code"))
			return
		}
		var codeInfo types.CodeInfo
		cliCtx.Cdc.MustUnmarshalBinaryBare(res, &codeInfo)
		rest.PostProcessResponseBare(w, cliCtx, codeInfo)
	}
}


func listContractsByCodeHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		accountAddr := r.FormValue("accountAddress")
		accountAddress := sdk.HexToAddress(accountAddr)
		params := types.NewContractListParams(accountAddress)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractList, bz)
		if err != nil {
			rest.WriteErrorRes(w, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "the length of contract list is 0"))
			return
		}
		var contractList types.ContractListResponse
		cliCtx.Cdc.MustUnmarshalBinaryBare(res, &contractList)

		rest.PostProcessResponseBare(w, cliCtx, contractList)
	}
}

func queryContractHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		contractAddr := r.FormValue("contractAddress")
		contractAddress := sdk.HexToAddress(contractAddr)

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok {
			rest.WriteErrorRes(w, err)
			return
		}

		params := types.NewQueryContractInfoParams(contractAddress)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractInfo, bz)
		if err != nil {
			rest.WriteErrorRes(w, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "no expected contract"))
			return
		}
		var contractInfo types.ContractInfo
		cliCtx.Cdc.MustUnmarshalBinaryBare(res, &contractInfo)
		rest.PostProcessResponseBare(w, cliCtx, contractInfo)
	}
}

func queryContractStateAllHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var param types.CallContractParam
		contractAddr := r.FormValue("contractAddress")
		contractAddress := sdk.HexToAddress(contractAddr)
		msg := r.FormValue("queryMSg")
		queryMsg, Err := hex.DecodeString(msg)
		if Err != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("get msg failed"))
			return
		}
		Err = json.Unmarshal(queryMsg, &param)
		if Err != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("unexpected query msg"))
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok {
			rest.WriteErrorRes(w, err)
			return
		}

		params := types.NewContractStateParam(contractAddress, queryMsg)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractState, bz)
		if err != nil {
			rest.WriteErrorRes(w, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var contractState types.ContractState
		cliCtx.Cdc.MustUnmarshalJSON(res, &contractState)
		rest.PostProcessResponseBare(w, cliCtx, contractState)
	}
}
