package rest

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	tRest "github.com/ci123chain/ci123chain/pkg/transfer/rest"
	sSdk "github.com/ci123chain/ci123chain/sdk/staking"
	"net/http"
	"strconv"
)

func RegisterRestTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/staking/validator/create", CreateValidatorRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/delegate", DelegateTX(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/redelegate", RedelegateTX(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/undelegate", UndelegateTX(cliCtx)).Methods("POST")
}

var cdc = app.MakeCodec()


func CreateValidatorRequest(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		gas, cAmt, msd, r, mr, mcr, moniker, identity, website, securityContact, details,
		priv, from, validatorAddr, delegatorAddr, publicKey, nonce, err := ParseArgs(request)

		txByte, err := sSdk.SignCreateValidatorMSg(from, cAmt, gas, nonce, priv, msd, validatorAddr,
			delegatorAddr, r, mr, mcr, moniker, identity, website, securityContact, details, publicKey)
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

func DelegateTX(cliCtx context.Context) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {

		key := request.FormValue("privateKey")
		from := request.FormValue("from")
		gas := request.FormValue("gas")
		Gas, err := strconv.ParseInt(gas, 10, 64)
		if err != nil || Gas < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"gas error"))
			return
		}

		isBalanceEnough := tRest.CheckBalanceFromParams(cliCtx, request)
		if !isBalanceEnough {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"The balance is not enough to pay the delegate"))
			return
		}

		UserGas := uint64(Gas)
		userNonce := request.FormValue("nonce")

		nonce, err := ParseNonce(from, userNonce)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"nonce error"))
			return
		}
		amount := request.FormValue("amount")
		amt, err := strconv.ParseInt(amount, 10, 64)
		if err != nil || amt < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"amount of coin error"))
			return
		}
		cAmt := uint64(amt)
		validatorAddr := request.FormValue("validatorAddr")
		delegatorAddr := request.FormValue("delegatorAddr")

		txByte, err := sSdk.SignDelegateMsg(from,cAmt, UserGas, nonce, key, validatorAddr, delegatorAddr)
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

func RedelegateTX(cliCtx context.Context) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {

		key := request.FormValue("privateKey")
		from := request.FormValue("from")
		gas := request.FormValue("gas")
		Gas, err := strconv.ParseInt(gas, 10, 64)
		if err != nil || Gas < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"gas error"))
			return
		}
		UserGas := uint64(Gas)
		userNonce := request.FormValue("nonce")
		nonce, err := ParseNonce(from, userNonce)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"nonce error"))
			return
		}
		amount := request.FormValue("amount")
		amt, err := strconv.ParseInt(amount, 10, 64)
		if err != nil || amt < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"amount of coin error"))
			return
		}
		cAmt := uint64(amt)
		validatorSrcAddr := request.FormValue("validatorSrcAddr")
		validatorDstAddr := request.FormValue("validatorDstAddr")
		delegatorAddr := request.FormValue("delegatorAddr")

		txByte, err := sSdk.SignRedelegateMsg(from,cAmt, UserGas, nonce, key, validatorSrcAddr, validatorDstAddr, delegatorAddr)
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

func UndelegateTX(cliCtx context.Context) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {

		key := request.FormValue("privateKey")
		from := request.FormValue("from")
		gas := request.FormValue("gas")
		Gas, err := strconv.ParseInt(gas, 10, 64)
		if err != nil || Gas < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"gas error"))
			return
		}
		UserGas := uint64(Gas)
		userNonce := request.FormValue("nonce")

		nonce, err := ParseNonce(from, userNonce)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"nonce error"))
			return
		}
		amount := request.FormValue("amount")
		amt, err := strconv.ParseInt(amount, 10, 64)
		if err != nil || amt < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"amount of coin error"))
			return
		}
		cAmt := uint64(amt)
		validatorAddr := request.FormValue("validatorAddr")
		delegatorAddr := request.FormValue("delegatorAddr")

		txByte, err := sSdk.SignUndelegateMsg(from,cAmt, UserGas, nonce, key, validatorAddr, delegatorAddr)
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

func ParseNonce(from, userNonce string) (uint64, error) {
	var nonce uint64
	froms, err := helper.StrToAddress(from)
	if err != nil {
		return nonce, err
	}
	if userNonce != "" {
		UserNonce, err := strconv.ParseInt(userNonce, 10, 64)
		if err != nil || UserNonce < 0 {
			return nonce, err
		}
		nonce = uint64(UserNonce)
		return nonce, nil
	}else {
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return nonce, err
		}
		nonce, err = ctx.GetNonceByAddress(froms)
		if err != nil {
			return nonce, err
		}
	}
	return nonce, nil
}

func ParseArgs(request *http.Request) (uint64,uint64, int64, int64, int64, int64, string,
	string, string, string, string, string, string, string, string, string, uint64, error) {
	gas := request.FormValue("gas")
	Gas, err := strconv.ParseInt(gas, 10, 64)
	if err != nil || Gas < 0 {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("gas error")
	}
	UserGas := uint64(Gas)
	amount := request.FormValue("amount")
	amt, err := strconv.ParseInt(amount, 10, 64)
	if err != nil || amt < 0 {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("amount error")
	}
	cAmt := uint64(amt)

	minSelfDelegation := request.FormValue("minSelfDelegation")
	msd, err := strconv.ParseInt(minSelfDelegation, 10, 64)
	if err != nil || msd < 0 {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("minSelfDelegation error")
	}
	rate := request.FormValue("rate")
	var r, mr, mcr int64
	if rate == "" {
		r = 1
	}else {
		r, err = strconv.ParseInt(rate, 10, 64)
		if err != nil || r < 0 {
			return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("rate error")
		}
	}
	maxRate := request.FormValue("maxRate")
	if maxRate == "" {
		mr = 1
	}else {
		mr, err = strconv.ParseInt(maxRate, 10, 64)
		if err != nil || mr < 0 {
			return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("max rate error")
		}
	}

	maxChangeRate := request.FormValue("maxChangeRate")
	if maxChangeRate == "" {
		mcr = 1
	}else {
		mcr, err = strconv.ParseInt(maxChangeRate, 10, 64)
		if err != nil || mcr < 0 {
			return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("max change rate error")
		}
	}
	moniker := request.FormValue("moniker")
	if moniker == "" {
		moniker = "moniker"
	}
	identity := request.FormValue("identity")
	if identity == "" {
		identity = "identity"
	}
	website := request.FormValue("website")
	if website == "" {
		website = "website"
	}
	securityContact := request.FormValue("securityContact")
	if securityContact == "" {
		securityContact = "securityContact"
	}
	details := request.FormValue("details")
	if details == "" {
		details = "details"
	}

	priv := request.FormValue("privateKey")
	if priv == "" {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("private key can't be empty")
	}
	from := request.FormValue("from")
	if from == "" {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("from error")
	}

	userNonce := request.FormValue("nonce")
	nonce, err := ParseNonce(from, userNonce)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("nonce error")
	}

	validatorAddr := request.FormValue("validatorAddress")
	if validatorAddr == "" {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("validator address error")
	}
	delegatorAddr := request.FormValue("delegatorAddress")
	if delegatorAddr == "" {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("delegator address error")
	}
	publicKey := request.FormValue("publicKey")
	if publicKey == "" {
		return 0, 0, 0, 0, 0, 0, "", "", "", "", "", "", "", "", "", "", 0, errors.New("public key can't be empty")
	}

	return UserGas, cAmt, msd, r, mr, mcr, moniker, identity, website, securityContact, details, priv, from, validatorAddr, delegatorAddr, publicKey, nonce, nil
}