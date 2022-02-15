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
	sSdk "github.com/ci123chain/ci123chain/sdk/staking"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func RegisterRestTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/preStaking/preDelegate", rest.MiddleHandler(cliCtx, PreDelegateRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/preStaking/delegator/delegate", rest.MiddleHandler(cliCtx, DelegateRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/preStaking/delegator/delegateDirect", rest.MiddleHandler(cliCtx, DelegateDirectRequest, types.DefaultCodespace)).Methods("POST")

	r.HandleFunc("/preStaking/delegator/redelegate", rest.MiddleHandler(cliCtx, RedelegateRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/preStaking/delegator/undelegate", rest.MiddleHandler(cliCtx, UndelegateRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/preStaking/create/validator", rest.MiddleHandler(cliCtx, CreateValidatorRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/preStaking/create/validatorDirect", rest.MiddleHandler(cliCtx, CreateValidatorDirectRequest, types.DefaultCodespace)).Methods("POST")

	r.HandleFunc("/preStaking/SetTokenAddress", rest.MiddleHandler(cliCtx, SetStakingTokenRequest, types.DefaultCodespace)).Methods("POST")

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

	dt := request.FormValue("delegate_time")
	t, err := time.ParseDuration(dt)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid delegate_time").Error())
		return
	}
	if t <= 0 {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid delegate_time").Error())
		return
	}

	msg := types.NewMsgPreStaking(from, coin, t)

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


func DelegateDirectRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
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
	dt := request.FormValue("delegate_time")
	t, err := time.ParseDuration(dt)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid delegate_time").Error())
		return
	}
	if t <= 0 {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid delegate_time").Error())
		return
	}

	v := request.FormValue("validator_address")
	validator := sdk.HexToAddress(v)

	msg := types.NewMsgStakingDirect(from, from, validator, coin, t)
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

func SetStakingTokenRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request)  {
	from := cliCtx.FromAddr
	err := checkAccountExist(cliCtx, cliCtx.FromAddr)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}
	tokenAddress := request.FormValue("token_address")
	if sdk.HexToAddress(tokenAddress).Empty() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token_address").Error())
		return
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, true)

	msg := types.NewMsgSetStakingToken(from, sdk.HexToAddress(tokenAddress))

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


func CreateValidatorDirectRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	delegatorAddr := cliCtx.FromAddr
	validatorAddr := cliCtx.FromAddr

	//verify account exists
	err := checkAccountExist(cliCtx, cliCtx.FromAddr)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}

	amount, ok := sdk.NewIntFromString(request.FormValue("amount"))
	if !ok {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	denom := request.FormValue("denom")
	coin := sdk.NewCoin(denom, amount)
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}

	dt := request.FormValue("delegate_time")
	t, err := time.ParseDuration(dt)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid delegate_time").Error())
		return
	}
	if t <= 0 {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid delegate_time").Error())
		return
	}


	msd, r, mr, mcr, moniker, identity, website, securityContact, details,
	publicKey,err := parseOtherArgs(request)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	MSD, R, MR, MXR := sSdk.CreateParseArgs(msd, r, mr, mcr)
	msg := types.NewMsgCreateValidatorDirect(cliCtx.FromAddr, MSD, validatorAddr,
		delegatorAddr, R, MR, MXR, moniker, identity, website, securityContact, details, publicKey, coin, t)

	if !cliCtx.Broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(cliCtx.FromAddr, cliCtx.Nonce, cliCtx.Gas, []sdk.Msg{msg}, cliCtx.PrivateKey, cdc)
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

func CreateValidatorRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	delegatorAddr := cliCtx.FromAddr
	validatorAddr := cliCtx.FromAddr

	//verify account exists
	err := checkAccountExist(cliCtx, cliCtx.FromAddr)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid from").Error())
		return
	}

	msd, r, mr, mcr, moniker, identity, website, securityContact, details,
	publicKey,err := parseOtherArgs(request)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	MSD, R, MR, MXR := sSdk.CreateParseArgs(msd, r, mr, mcr)
	id := request.FormValue("vault_id")
	msg := types.NewMsgCreateValidator(cliCtx.FromAddr, MSD, validatorAddr,
		delegatorAddr, R, MR, MXR, moniker, identity, website, securityContact, details, publicKey, id)

	if !cliCtx.Broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(cliCtx.FromAddr, cliCtx.Nonce, cliCtx.Gas, []sdk.Msg{msg}, cliCtx.PrivateKey, cdc)
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

func parseOtherArgs(request *http.Request) (int64, int64, int64, int64, string,
	string, string, string, string, string, error) {

	minSelfDelegation := request.FormValue("min_self_delegation")
	msd, err := strconv.ParseInt(minSelfDelegation, 10, 64)
	if err != nil || msd < 0 {
		return 0, 0, 0, 0, "", "", "", "", "", "",errors.New("min_self_delegation error")
	}
	rate := request.FormValue("rate")
	var r, mr, mcr int64
	if rate == "" {
		r = 1
	}else {
		r, err = strconv.ParseInt(rate, 10, 64)
		if err != nil || r < 0 {
			return 0, 0, 0, 0, "", "", "", "", "", "", errors.New("rate error")
		}
	}
	maxRate := request.FormValue("max_rate")
	if maxRate == "" {
		mr = 1
	}else {
		mr, err = strconv.ParseInt(maxRate, 10, 64)
		if err != nil || mr < 0 {
			return 0, 0, 0, 0, "", "", "", "", "", "",  errors.New("max_rate error")
		}
	}

	maxChangeRate := request.FormValue("max_change_rate")
	if maxChangeRate == "" {
		mcr = 1
	}else {
		mcr, err = strconv.ParseInt(maxChangeRate, 10, 64)
		if err != nil || mcr < 0 {
			return 0, 0, 0, 0, "", "", "", "", "", "",errors.New("max_change_rate error")
		}
	}
	moniker := request.FormValue("moniker")
	identity := request.FormValue("identity")
	website := request.FormValue("website")
	securityContact := request.FormValue("security_contact")
	details := request.FormValue("details")

	publicKey := request.FormValue("public_key")
	if publicKey == "" {
		return 0, 0, 0, 0, "", "", "", "", "", "", errors.New("public_key can't be empty")
	}

	return msd, r, mr, mcr, moniker, identity, website, securityContact, details, publicKey, nil
}
