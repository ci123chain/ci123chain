package rest

import (
	"encoding/hex"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"net/http"
)

func BroadcastTxRequest(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		data := request.FormValue("tx_byte")
		err := util.CheckStringLength(1, 1000, data)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"error txByte length"))
			return
		}
		txByte, err := hex.DecodeString(data)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"error txByte"))
			return
		}

		res, err := cliCtx.BroadcastSignedData(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
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
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
			return
		}

		_, _ = cliCtx.BroadcastSignedDataAsync(txByte)
		/*
		res, err := cliCtx.BroadcastSignedDataAsync(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, clients.ErrBroadcast(types.DefaultCodespace, err))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
		*/
	}
}
