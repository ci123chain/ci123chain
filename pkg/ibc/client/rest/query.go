package rest

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"io/ioutil"
	"net/http"
)


func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/ibctx", QueryTxByUniqueIDRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/ibctx/state", QueryTxByStateRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/ibctx/nonce", QueryAccountNonceRequestHandlerFn(cliCtx)).Methods("POST")

}

type Params struct {
	Data string `json:"data"`
}

type IBCUniqueIDData struct {
	UniqueID	string	`json:"uniqueID"`
}
func QueryTxByStateRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		//ibcState := vars["ibcstate"]
		var params Params
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &params)
		if readErr != nil {
			//
		}

		if err := types.ValidateState(params.Data); err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request)
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/state/" + params.Data, nil)
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "There is no ibctx ready"))
			return
		}
		var ibcMsg types.IBCInfo
		err2 := cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(res, &ibcMsg)
		if err2 != nil {
			rest.WriteErrorRes(writer, transfer.ErrCheckParams(types.DefaultCodespace, err2.Error()))
			return
		}
		resp := &IBCUniqueIDData{UniqueID:string(ibcMsg.UniqueID)}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

type IBCTxStateData struct {
	State	string	`json:"state"`
}
func QueryTxByUniqueIDRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		//uniqueidStr := vars["uniqueid"]

		var params Params
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &params)
		if readErr != nil {
			//
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request)
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}
		uniqueBz := []byte(params.Data)

		res, _, err := cliCtx.Query("/store/" + types.StoreKey + "/types", uniqueBz)
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "this uniqueID is not exist"))
			return
		}
		var ibcMsg types.IBCInfo
		err2 := cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(res, &ibcMsg)
		if err2 != nil {
			rest.WriteErrorRes(writer, transfer.ErrCheckParams(types.DefaultCodespace, err2.Error()))
			return
		}
		if !bytes.Equal(uniqueBz, ibcMsg.UniqueID) {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, fmt.Sprintf("different uniqueID get %s, expected %s", hex.EncodeToString(ibcMsg.UniqueID), params.Data)))
			return
		}
		resp := &IBCTxStateData{State:ibcMsg.State}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

type NonceData struct {
	Nonce 	uint64 `json:"nonce"`
}
func QueryAccountNonceRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		//accountAddress := vars["accountaddress"]
		var params Params
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &params)
		if readErr != nil {
			//
		}
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request)
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/nonce/" + params.Data, nil)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var nonce uint64
		err2 := cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(res, &nonce)
		if err2 != nil {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, err2.Error()))
			return
		}
		resp := &NonceData{Nonce:nonce}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}
