package rest

import (
	//"encoding/hex"
	"github.com/pkg/errors"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/util"
	"strconv"

	///sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
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
		isBalanceEnough := CheckBalanceFromParams(cliCtx, request)
		if !isBalanceEnough {
			rest.WriteErrorRes(writer, transaction.ErrAmount(types.DefaultCodespace, errors.New("The balance is not enough to pay the amount")) )
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

// check balance from tranfer params

func CheckBalanceFromParams(ctx context.Context, r *http.Request) bool {
	from := r.FormValue("from")
	amount := r.FormValue("amount")

	acc, _ := helper.StrToAddress(from)
	balance, _, err := ctx.GetBalanceByAddress(acc)
	if err != nil {
		return false
	}
	amountU, _ := strconv.ParseUint(amount,10,64)
	if balance < amountU {
		return false
	}
	return true

}