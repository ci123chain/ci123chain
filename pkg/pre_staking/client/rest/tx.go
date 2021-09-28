package rest

import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func RegisterRestTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/preStaking/preDelegate", rest.MiddleHandler(cliCtx, PreDelegateRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/preStaking/delegator/delegate", rest.MiddleHandler(cliCtx, DelegateRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/preStaking/delegator/redelegate", rest.MiddleHandler(cliCtx, RedelegateRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/preStaking/delegator/undelegate", rest.MiddleHandler(cliCtx, UndelegateRequest, types.DefaultCodespace)).Methods("POST")
	//r.HandleFunc("/preStaking/deploy", rest.MiddleHandler(cliCtx, DeployRequest, types.DefaultCodespace)).Methods("POST")
}


var cdc = types2.GetCodec()

func PreDelegateRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	amount, ok := sdk.NewIntFromString(request.FormValue("amount"))
	if !ok {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}

	denom := request.FormValue("denom")
	coin := sdk.NewCoin(denom, amount)
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	c := request.FormValue("contract")
	contract := sdk.HexToAddress(c)

	dt := request.FormValue("delegate_time")
	t, err := time.ParseDuration(dt)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}
	if t <= 0 {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}

	msg := types.NewMsgPreStaking(from, coin, contract, t)

	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}


func DelegateRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	//amount, ok := sdk.NewIntFromString(request.FormValue("amount"))
	//if !ok {
	//	rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
	//	return
	//}
	////verify account exists
	//err = checkAccountExist(cliCtx, from)
	//if err != nil {
	//	rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
	//	return
	//}
	//
	//denom := request.FormValue("denom")
	//coin := sdk.NewCoin(denom, amount)
	//if coin.IsNegative() || coin.IsZero() {
	//	rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
	//	return
	//}
	id := request.FormValue("vault_id")
	//Id, ok := new(big.Int).SetString(id, 10)
	//if !ok {
	//	rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid vault_id").Error())
	//	return
	//}
	v := request.FormValue("validator_address")
	validator := sdk.HexToAddress(v)

	msg := types.NewMsgStaking(from, from, validator, id)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}


func RedelegateRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}
	src := request.FormValue("src_validator_address")
	srcValidator := sdk.HexToAddress(src)
	dst := request.FormValue("dst_validator_address")
	dstValidator := sdk.HexToAddress(dst)


	msg := types.NewMsgRedelegate(from, srcValidator, dstValidator)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}


func UndelegateRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	//amount, ok := sdk.NewIntFromString(request.FormValue("amount"))
	//if !ok {
	//	rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
	//	return
	//}
	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}

	//denom := request.FormValue("denom")
	//coin := sdk.NewCoin(denom, amount)
	//if coin.IsNegative() || coin.IsZero() {
	//	rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
	//	return
	//}
	id := request.FormValue("vault_id")
	//Id, ok := new(big.Int).SetString(id, 10)
	//if !ok {
	//	rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid vault_id").Error())
	//	return
	//}
	msg := types.NewMsgUndelegate(from, id)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func DeployRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}

	msg := types.NewMsgDeploy(from)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func checkAccountExist(ctx context.Context, address... sdk.AccAddress) error {
	for i := 0; i < len(address); i++ {
		_,_,  err := ctx.GetNonceByAddress(address[i], false)
		if err != nil {
			return errors.New(fmt.Sprintf("account of %s does not exist", address[i]))
		}
	}
	return nil
}