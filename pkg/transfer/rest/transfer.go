package rest

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	transfer2 "github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	fabric := request.FormValue("fabric")
	isFabric, err := util.CheckFabric(fabric)
	if err != nil {
		isFabric = false
	}
	denom := request.FormValue("denom")
	if denom == "" {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "the name of the coin which you want to transfer cant not be empty").Error())
		return
	}
	to := sdk.HexToAddress(request.FormValue("to"))
	amount, err := strconv.ParseUint(request.FormValue("amount"), 10, 64)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	isBalanceEnough := CheckAccountAndBalanceFromParams(cliCtx, request, writer, denom)
	if !isBalanceEnough {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "The balance is not enough to pay the amount").Error())
		return
	}
	coin := sdk.NewUInt64Coin(denom, amount)
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	msg := transfer2.NewMsgTransfer(from, to, sdk.NewCoins(coin), isFabric)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	var a types2.CommonTx
	err = cdc.UnmarshalBinaryBare(txByte, &a)
	if err != nil {
		fmt.Println(err)
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "broadcast failed").Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

// check balance from tranfer params

func CheckAccountAndBalanceFromParams(ctx context.Context, r *http.Request, w http.ResponseWriter, denom string) bool {
	from := r.FormValue("from")
	amount := r.FormValue("amount")


	acc, _ := helper.StrToAddress(from)
	balance, _, err := ctx.GetBalanceByAddress(acc, false, "")
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, "get balances of from account failed").Error())
		return false
	}
	amountI, ok := sdk.NewIntFromString(amount)
	if !ok {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return false
	}
	if balance.AmountOf(denom).LT(amountI) {
		return false
	}
	return true
}