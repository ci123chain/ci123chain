package rest

import (
	"encoding/hex"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transfer/types"
	"net/http"
)

func BroadcastTxRequest(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data := request.FormValue("data")
		txByte, err := hex.DecodeString(data)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"data error"))
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
		data := request.FormValue("data")
		txByte, err := hex.DecodeString(data)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"data error"))
			return
		}

		res, err := cliCtx.BroadcastSignedDataAsync(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
	}
}
