package rest

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/ibc/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/gorilla/mux"
	"net/http"
)


func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/ibctx", QueryTxByUniqueIDRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/ibctx/state", QueryTxByStateRequestHandlerFn(cliCtx)).Methods("POST")

}

type IBCUniqueIDData struct {
	UniqueID	string	`json:"uniqueID"`
}

type StateParams struct {
	State     string `json:"state"`
	Height    string `json:"height"`
}

type QueryStateParams struct {
	//
	Data StateParams `json:"data"`
}
func QueryTxByStateRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ibcState := request.FormValue("ibcstate")
		height := request.FormValue("height")
		prove := request.FormValue("prove")
		if err := types.ValidateState(ibcState); err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}

		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/state/" + ibcState, nil, isProve)
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
		value := &IBCUniqueIDData{UniqueID:string(ibcMsg.UniqueID)}
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

type IBCTxStateData struct {
	State	string	`json:"state"`
}

type Txparams struct {
	UniqueID   string  `json:"unique_id"`
	Height     string  `json:"height"`
}

type QueryTxParams struct {
	Data Txparams `json:"data"`
}

func QueryTxByUniqueIDRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		UniqueidStr := request.FormValue("uniqueID")
		height := request.FormValue("height")
		prove := request.FormValue("prove")
		checkErr := util.CheckStringLength(1, 100, UniqueidStr)
		if checkErr != nil {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "unexpected uniqueID"))
			return
		}
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}
		uniqueBz := []byte(UniqueidStr)

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/store/" + types.StoreKey + "/types", uniqueBz, isProve)
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
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, fmt.Sprintf("different uniqueID get %s, expected %s", hex.EncodeToString(ibcMsg.UniqueID), UniqueidStr)))
			return
		}
		value := &IBCTxStateData{State:ibcMsg.State}
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}
