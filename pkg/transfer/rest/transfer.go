package rest

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	transfer2 "github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func SendRequestHandlerFn(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	fabric := request.FormValue("fabric")
	isFabric, err := util.CheckFabric(fabric)
	if err != nil {
		isFabric = false
	}
	to := sdk.HexToAddress(request.FormValue("to"))
	amount, err := strconv.ParseUint(request.FormValue("amount"), 10, 64)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	isBalanceEnough, denom := CheckAccountAndBalanceFromParams(cliCtx, request, writer)
	if !isBalanceEnough {
		rest.WriteErrorRes(writer, transaction.ErrAmount(types.DefaultCodespace, errors.New("The balance is not enough to pay the amount")) )
		return
	}
	coin := sdk.NewUInt64Coin(denom, amount)
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, "invalid amount"))
		return
	}
	msg := transfer2.NewMsgTransfer(from, to, coin, isFabric)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	var a types2.CommonTx
	err = cdc.UnmarshalBinaryBare(txByte, &a)
	if err != nil {
		fmt.Println(err)
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

// check balance from tranfer params

func CheckAccountAndBalanceFromParams(ctx context.Context, r *http.Request, w http.ResponseWriter) (bool, string) {
	from := r.FormValue("from")
	amount := r.FormValue("amount")


	acc, _ := helper.StrToAddress(from)
	balance, _, err := ctx.GetBalanceByAddress(acc, false)

	if err != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("get balances of from account failed"))
		return false, ""
	}
	amountI, ok := sdk.NewIntFromString(amount)
	if !ok {
		rest.WriteErrorRes(w, sdk.ErrInternal(fmt.Sprintf("invalid amount %s", amount)))
		return false, ""
	}
	if balance.Amount.LT(amountI) {
		return false, ""
	}
	return true, balance.Denom

}