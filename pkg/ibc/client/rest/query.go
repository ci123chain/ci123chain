package rest

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"net/http"
)


func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/ibctx/{uniqueid}", QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
}

func QueryTxRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
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
