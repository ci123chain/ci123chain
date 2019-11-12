package rest

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/ibc/keeper"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"net/http"
	"strings"
)


func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/ibctx/{uniqueid}", QueryTxByUniqueIDRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/ibctx/state/{ibcstate}", QueryTxByStateRequestHandlerFn(cliCtx)).Methods("GET")

}

func QueryTxByStateRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		ibcState := vars["ibcstate"]

		if err := keeper.ValidateState(ibcState); err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request)
		if !ok {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, "Build context error")
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/state/" + ibcState, nil)
		if len(res) < 1 {
			rest.WriteErrorResponse(writer, http.StatusNotFound, "no ready ibc tx found " )
		}
		var ibcMsg types.IBCMsg
		err = cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(res, &ibcMsg)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
		}
		rest.PostProcessResponseBare(writer, cliCtx, strings.ToUpper(hex.EncodeToString(ibcMsg.UniqueID)))
	}
}



func QueryTxByUniqueIDRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		uniqueidStr := vars["uniqueid"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request)

		if !ok {
			return
		}
		uniqueBz , err := hex.DecodeString(uniqueidStr)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
		}
		res, _, err := cliCtx.Query("/store/" + types.StoreKey + "/types", uniqueBz)
		if len(res) < 1 {
			rest.WriteErrorResponse(writer, http.StatusNotFound, "no ibc tx found with uniqueid " + uniqueidStr)
		}
		var ibcMsg types.IBCMsg
		err = cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(res, &ibcMsg)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
		}
		if !bytes.Equal(uniqueBz, ibcMsg.UniqueID) {
			rest.WriteErrorResponse(writer, http.StatusNotFound, fmt.Sprintf("different uniqueid get %s, expected %s", hex.EncodeToString(ibcMsg.UniqueID), uniqueidStr))
		}
		rest.PostProcessResponseBare(writer, cliCtx, ibcMsg.State)
	}
}
