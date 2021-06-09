package rest

import (
	"encoding/hex"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/util"
	"net/http"
)

func BroadcastTxRequest(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		data := request.FormValue("tx_byte")
		err := util.CheckStringLength(1, 1000, data)
		if err != nil {
			rest.WriteErrorRes(writer, "invalid tx_byte")
			return
		}
		txByte, err := hex.DecodeString(data)
		if err != nil {
			rest.WriteErrorRes(writer, "invalid tx_byte")
			return
		}

		res, err := cliCtx.BroadcastSignedData(txByte)
		if err != nil {
			rest.WriteErrorRes(writer,err.Error())
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
	}
}

func BroadcastTxRequestAsync(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		data := request.FormValue("tx_byte")
		txByte, err := hex.DecodeString(data)
		if err != nil {
			rest.WriteErrorRes(writer, "invalid tx_byte")
			return
		}


		res, err := cliCtx.BroadcastSignedDataAsync(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
		
	}
}
