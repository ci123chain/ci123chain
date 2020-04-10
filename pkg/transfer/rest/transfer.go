package rest

import (
	//"encoding/hex"
	"github.com/pkg/errors"
	"github.com/tanhuiya/ci123chain/pkg/util"

	///sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/transfer/types"
	"net/http"
)



func SendRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		priv := request.FormValue("privateKey")
		err := util.CheckStringLength(1, 100, priv)
		if err != nil {
			rest.WriteErrorRes(writer, transaction.ErrBadPrivkey(types.DefaultCodespace, errors.New("param privateKey not found")) )
			return
		}
		async := request.FormValue("async")
		ok, err := util.CheckBool(async)  //default async
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"error async"))
			return
		}
		txByte, err := buildTransferTx(request, false, priv)
		if err != nil {
			rest.WriteErrorRes(writer, transaction.ErrSignature(types.DefaultCodespace, errors.New("sign with tx error")))
			return
		}

		/*res, err := cliCtx.BroadcastSignedData(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)*/

		if ok {
			//async
			res, err := cliCtx.BroadcastTxAsync(txByte)
			if err != nil {
				rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
				return
			}
			rest.PostProcessResponseBare(writer, cliCtx, res)
		}else {
			//sync
			res, err := cliCtx.BroadcastSignedData(txByte)
			if err != nil {
				rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
				return
			}
			rest.PostProcessResponseBare(writer, cliCtx, res)
		}
	}
}