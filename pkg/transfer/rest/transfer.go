package rest

import (
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/util"
	//"encoding/hex"
	"github.com/pkg/errors"
	"strconv"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	///sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	"net/http"
)

func SendRequestHandlerFn(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	priv := request.FormValue("privateKey")
	err := util.CheckStringLength(1, 100, priv)
	if err != nil {
		rest.WriteErrorRes(writer, transaction.ErrBadPrivkey(types.DefaultCodespace, errors.New("param privateKey not found")) )
		return
	}

	txByte, err := buildTransferTx(request, false, priv)
	if err != nil {
		rest.WriteErrorRes(writer, transaction.ErrSignature(types.DefaultCodespace, err))
		return
	}
	isBalanceEnough := CheckAccountAndBalanceFromParams(cliCtx, request, writer)
	if !isBalanceEnough {
		rest.WriteErrorRes(writer, transaction.ErrAmount(types.DefaultCodespace, errors.New("The balance is not enough to pay the amount")) )
		return
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return

	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

// check balance from tranfer params

func CheckAccountAndBalanceFromParams(ctx context.Context, r *http.Request, w http.ResponseWriter) bool {
	from := r.FormValue("from")
	amount := r.FormValue("amount")


	acc, _ := helper.StrToAddress(from)
	balance, _, err := ctx.GetBalanceByAddress(acc, false)

	if err != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("get balances of from account failed"))
		return false
	}
	amountU, _ := strconv.ParseUint(amount,10,64)
	if balance < amountU {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"The balance is not enough to pay the delegate"))
		return false
	}
	return true

}