package rest

import (
	"CI123Chain/pkg/abci/types/rest"
	"CI123Chain/pkg/client/context"
	"io/ioutil"
	"net/http"
)


// BroadcastReq defines a tx broadcasting request.
type BroadcastReq struct {
	//Tx   types.StdTx `json:"tx" yaml:"tx"`
	Mode string      `json:"mode" yaml:"mode"`
}


func BroadcastTxRequest(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Cdc.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// todo req.Tx
		txBytes, err := cliCtx.Cdc.MarshalBinaryLengthPrefixed(req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}
