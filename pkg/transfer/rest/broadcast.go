package rest

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"encoding/hex"
	"net/http"
)

func BraodcastTxRequest(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data := request.FormValue("data")
		txByte, err := hex.DecodeString(data)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, "invalid data")
			return
		}

		res, err := cliCtx.BroadcastSignedData(txByte)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
	}
}
