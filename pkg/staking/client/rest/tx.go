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
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	sSdk "github.com/ci123chain/ci123chain/sdk/staking"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func RegisterRestTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/staking/validator/create", rest.MiddleHandler(cliCtx, CreateValidatorRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/staking/delegator/delegate", rest.MiddleHandler(cliCtx, DelegateTX, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/staking/delegator/redelegate", rest.MiddleHandler(cliCtx, RedelegateTX, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/staking/delegator/undelegate", rest.MiddleHandler(cliCtx, UndelegateTX, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/staking/validator/edit", rest.MiddleHandler(cliCtx, EditValidatorTX, types.DefaultCodespace)).Methods("POST")
}

var cdc = types2.MakeCodec()

func CreateValidatorRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	delegatorAddr := from
	validatorAddr := from
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

	msd, r, mr, mcr, moniker, identity, website, securityContact, details,
	publicKey,err := parseOtherArgs(request)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	coin := sdk.NewChainCoin(amount)
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	MSD, R, MR, MXR := sSdk.CreateParseArgs(msd, r, mr, mcr)
	msg := staking.NewCreateValidatorMsg(from, coin, MSD, validatorAddr,
		delegatorAddr, R, MR, MXR, moniker, identity, website, securityContact, details, publicKey)

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

func DelegateTX(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	delegatorAddr := from
	validatorAddr := sdk.HexToAddress(request.FormValue("validator_address"))
	amount, ok := sdk.NewIntFromString(request.FormValue("amount"))
	if !ok {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	coin := sdk.NewChainCoin(amount)
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from, validatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error()).Error())
		return
	}


	msg := staking.NewDelegateMsg(from, delegatorAddr, validatorAddr, coin)
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

func RedelegateTX(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	delegatorAddr := from
	validatorSrcAddr := sdk.HexToAddress(request.FormValue("validator_src_address"))
	validatorDstAddr := sdk.HexToAddress(request.FormValue("validator_dst_address"))
	amount, ok := sdk.NewIntFromString(request.FormValue("amount"))
	if !ok {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from, validatorSrcAddr, validatorDstAddr)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error()).Error())
		return
	}

	coin := sdk.NewChainCoin(amount)
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	msg := staking.NewRedelegateMsg(from, delegatorAddr, validatorSrcAddr, validatorDstAddr, coin)

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
		rest.WriteErrorRes(writer,sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func UndelegateTX(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	delegatorAddr := from
	validatorAddr := sdk.HexToAddress(request.FormValue("validator_address"))
	amount, ok := sdk.NewIntFromString(request.FormValue("amount"))
	if !ok {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from, validatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error()).Error())
		return
	}

	coin := sdk.NewChainCoin(amount)
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}
	msg := staking.NewUndelegateMsg(from, delegatorAddr, validatorAddr, coin)
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

func EditValidatorTX(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error()).Error())
		return
	}
	moniker, identity, website, secu, details := getDescription(request)
	minSelf, newRate, ok := getMinSelfAndNewRate(writer, request)
	if !ok {
		return
	}
	var nrArg *sdk.Dec
	var minArg *sdk.Int
	if newRate < 0 {
		if newRate == -1 {
			nrArg = nil
		}else {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid new_rate").Error())
			return
		}
	}else {
		nr := sdk.NewDecWithPrec(newRate, 2)
		nrArg = &nr
	}
	if minSelf < 0 {
		if minSelf == -1 {
			minArg = nil
		}else {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "min_self_delegation error").Error())
			return
		}
	}else {
		min := sdk.NewInt(minSelf)
		minArg = &min
	}
	desc := types.Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: secu,
		Details:         details,
	}
	msg := staking.NewEditValidatorMsg(from, desc, nrArg, minArg)
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

func getMinSelfAndNewRate(w http.ResponseWriter, req *http.Request) (int64, int64, bool) {
	//
	var minSelf, newRate int64
	var err error
	ms := req.FormValue("min_self_delegation")
	if ms != "" {
		minSelf, err = util.CheckInt64(ms)
		if err != nil || minSelf <= 0 {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid min_self_delegation").Error())
			return 0, 0, false
		}
	}else {
		minSelf = -1
	}
	nr := req.FormValue("new_rate")
	if nr != "" {
		newRate, err = util.CheckInt64(nr)
		if err != nil || newRate > 100 || newRate <= 0 {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid new_rate").Error())
			return 0, 0, false
		}
	}else {
		newRate = -1
	}
	return minSelf, newRate, true
}

func getDescription(req *http.Request) (string, string, string, string, string) {
	moniker := req.FormValue("moniker")
	if moniker == "" {
		moniker = types.DoNotModifyDesc
	}
	identity := req.FormValue("identity")
	if identity == "" {
		identity = types.DoNotModifyDesc
	}
	website := req.FormValue("website")
	if website == "" {
		website = types.DoNotModifyDesc
	}
	secu := req.FormValue("security_contact")
	if secu == "" {
		secu = types.DoNotModifyDesc
	}
	details := req.FormValue("details")
	if details == "" {
		details = types.DoNotModifyDesc
	}
	return moniker, identity, website, secu, details
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

