package rest

import (
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	tRest "github.com/ci123chain/ci123chain/pkg/transfer/rest"
	"github.com/ci123chain/ci123chain/pkg/util"
	sSdk "github.com/ci123chain/ci123chain/sdk/staking"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func RegisterRestTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/staking/validator/create", rest.MiddleHandler(cliCtx, CreateValidatorRequest, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/staking/delegate", rest.MiddleHandler(cliCtx, DelegateTX, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/staking/redelegate", rest.MiddleHandler(cliCtx, RedelegateTX, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/staking/undelegate", rest.MiddleHandler(cliCtx, UndelegateTX, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/staking/edit", rest.MiddleHandler(cliCtx, EditValidatorTX, types.DefaultCodespace)).Methods("POST")

}

var cdc = app.MakeCodec()


func CreateValidatorRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	/*
	gas, cAmt, msd, r, mr, mcr, moniker, identity, website, securityContact, details,
	priv, from, validatorAddr, delegatorAddr, publicKey, nonce, err := ParseArgs(request)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}

	txByte, err := sSdk.SignCreateValidatorMSg(from, cAmt, gas, nonce, priv, msd, validatorAddr,
		delegatorAddr, r, mr, mcr, moniker, identity, website, securityContact, details, publicKey)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	*/
	msd, r, mr, mcr, moniker, identity, website, securityContact, details,
	publicKey,err := parseOtherArgs(request)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}

	from, delegatorAddr, key, amount, gas, nonce, ok := parseBaseArgs(writer, request)
	if !ok {
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
	}


	txByte, err := sSdk.SignCreateValidatorMSg(from, amount, gas, nonce, key, msd, from,
		delegatorAddr, r, mr, mcr, moniker, identity, website, securityContact, details, publicKey)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"data error"))
		return
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func DelegateTX(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	/*
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
	delegatorAddr := from//request.FormValue("delegatorAddr")


	txByte, err := sSdk.SignDelegateMsg(from,cAmt, UserGas, nonce, key, validatorAddr, delegatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	*/
	isOk := tRest.CheckAccountAndBalanceFromParams(cliCtx, request, writer)
	if !isOk {
		return
	}

	from, delegatorAddr, key, amount, gas, nonce, ok := parseBaseArgs(writer, request)
	if !ok {
		return
	}
	validatorAddr := request.FormValue("validatorAddress")
	//verify account exists
	err := checkAccountExist(cliCtx, from, validatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
	}

	txByte, err := sSdk.SignDelegateMsg(from, amount, gas, nonce, key, validatorAddr, delegatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"data error"))
		return
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func RedelegateTX(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	/*key := request.FormValue("privateKey")
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
	delegatorAddr := from///request.FormValue("delegatorAddr")

	txByte, err := sSdk.SignRedelegateMsg(from,cAmt, UserGas, nonce, key, validatorSrcAddr, validatorDstAddr, delegatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}*/
	validatorSrcAddr := request.FormValue("validatorSrcAddr")
	validatorDstAddr := request.FormValue("validatorDstAddr")

	from, delegatorAddr, key, amount, gas, nonce, ok := parseBaseArgs(writer, request)
	if !ok {
		return
	}
	//verify account exists
	err := checkAccountExist(cliCtx, from, validatorSrcAddr, validatorDstAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
	}


	txByte, err := sSdk.SignRedelegateMsg(from, amount, gas, nonce, key, validatorSrcAddr, validatorDstAddr, delegatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"data error"))
		return
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func UndelegateTX(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	/*
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
	delegatorAddr := from//request.FormValue("delegatorAddr")

	txByte, err := sSdk.SignUndelegateMsg(from,cAmt, UserGas, nonce, key, validatorAddr, delegatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	*/
	from, delegatorAddr, key, amount, gas, nonce, ok := parseBaseArgs(writer, request)
	if !ok {
		return
	}
	validatorAddr := request.FormValue("validatorAddress")
	//verify account exists
	err := checkAccountExist(cliCtx, from, validatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
	}

	txByte, err := sSdk.SignUndelegateMsg(from, amount, gas, nonce, key, validatorAddr, delegatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"data error"))
		return
	}


	/*txByte, err := sSdk.SignRedelegateMsg(from, amount, gas, nonce, key, validatorSrcAddr, validatorDstAddr, delegatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"data error"))
		return
	}*/

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}


func EditValidatorTX(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	from, gas, nonce, priv, ok := parseBasicArgs(writer, request)
	if !ok {
		return
	}
	//verify account exists
	err := checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	moniker, identity, website, secu, details := getDescription(request)
	minSelf, newRate, ok := getMinSelfAndNewRate(writer, request)
	if !ok {
		return
	}

	txByte, err := sSdk.SignEditValidator(from,gas, nonce, priv, moniker, identity, website, secu, details, minSelf, newRate)
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

func getMinSelfAndNewRate(w http.ResponseWriter, req *http.Request) (int64, int64, bool) {
	//
	var minSelf, newRate int64
	var err error
	ms := req.FormValue("minSelfDelegation")
	if ms != "" {
		minSelf, err = util.CheckInt64(ms)
		if err != nil || minSelf <= 0 {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid minSelfDelegation"))
			return 0, 0, false
		}
	}else {
		minSelf = -1
	}
	nr := req.FormValue("newRate")
	if nr != "" {
		newRate, err = util.CheckInt64(nr)
		if err != nil || newRate > 100 || newRate <= 0 {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid newRate"))
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
	secu := req.FormValue("securityContact")
	if secu == "" {
		secu = types.DoNotModifyDesc
	}
	details := req.FormValue("details")
	if details == "" {
		details = types.DoNotModifyDesc
	}
	return moniker, identity, website, secu, details
}

func parseBasicArgs(w http.ResponseWriter, req *http.Request) (string, uint64, uint64, string, bool) {
	/*from, ok := util.CheckFromAddressVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid from address"))
		return "", 0, 0, "", false
	}
	gas, ok := util.CheckGasVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid gas"))
		return "", 0, 0, "", false
	}
	userNonce := req.FormValue("nonce")
	nonce, err := ParseNonce(from, userNonce)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,"nonce error"))
		return "", 0, 0, "", false
	}
	privateKey, ok := util.CheckPrivateKey(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid private key"))
		return "", 0, 0, "", false
	}*/
	from, ok := util.CheckFromAddressVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid from address"))
		return "", 0, 0, "", false
	}
	gas, ok := util.CheckGasVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid gas"))
		return "", 0, 0, "", false
	}
	nonce, ok := checkNonce(req, sdk.HexToAddress(from), cdc)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid nonce"))
		return "", 0, 0, "", false
	}
	privateKey, ok := util.CheckPrivateKey(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid private key"))
		return "", 0, 0, "", false
	}
	return from, gas, nonce, privateKey, true
}

func checkNonce(r *http.Request, from sdk.AccAddress, cdc *codec.Codec) (uint64, bool) {
	nonce := r.FormValue("nonce")
	var Nonce uint64
	if nonce == "" {
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return 0, false
		}
		var Err error
		Nonce, _, Err = ctx.GetNonceByAddress(from, false)
		if Err != nil {
			return 0, false
		}
	}else {
		var checkErr error
		Nonce, checkErr = util.CheckUint64(nonce)
		if checkErr != nil {
			return 0, false
		}
	}
	return Nonce, true
}


func checkAccountExist(ctx context.Context, address... string) error {
	for i := 0; i < len(address); i++ {
		_,_,  err := ctx.GetNonceByAddress(sdk.HexToAddress(address[i]), false)
		if err != nil {
			return errors.New(fmt.Sprintf("account of %s does not exist", address[i]))
		}
	}
	return nil
}



func parseOtherArgs(request *http.Request) (int64, int64, int64, int64, string,
	string, string, string, string, string, error) {

	minSelfDelegation := request.FormValue("minSelfDelegation")
	msd, err := strconv.ParseInt(minSelfDelegation, 10, 64)
	if err != nil || msd < 0 {
		return 0, 0, 0, 0, "", "", "", "", "", "",errors.New("minSelfDelegation error")
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
	maxRate := request.FormValue("maxRate")
	if maxRate == "" {
		mr = 1
	}else {
		mr, err = strconv.ParseInt(maxRate, 10, 64)
		if err != nil || mr < 0 {
			return 0, 0, 0, 0, "", "", "", "", "", "",  errors.New("max rate error")
		}
	}

	maxChangeRate := request.FormValue("maxChangeRate")
	if maxChangeRate == "" {
		mcr = 1
	}else {
		mcr, err = strconv.ParseInt(maxChangeRate, 10, 64)
		if err != nil || mcr < 0 {
			return 0, 0, 0, 0, "", "", "", "", "", "",errors.New("max change rate error")
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

	publicKey := request.FormValue("publicKey")
	if publicKey == "" {
		return 0, 0, 0, 0, "", "", "", "", "", "", errors.New("public key can't be empty")
	}

	return msd, r, mr, mcr, moniker, identity, website, securityContact, details, publicKey, nil
}




func parseBaseArgs(w http.ResponseWriter, req *http.Request) (string, string, string, int64, uint64, uint64, bool) {

	//为了确保 createValidator交易是from账户发出的，validator的地址就直接是这个from;
	//只能是 设置 提取奖金的地址；那么就要保证 from validator delegator是一致的；但是 又没有验证；因为它本身就是使用相同的值；
	//要保证delegator跟from一致；
	/*from, ok := util.CheckFromAddressVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid from address"))
		return "", "", "", 0, 0, 0, false
	}*/
	from, gas, nonce, privateKey, ok := parseBasicArgs(w, req)
	if !ok {
		return "", "", "", 0, 0, 0, false
	}
	amount, ok := util.CheckAmountVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid amount"))
		return "", "", "", 0, 0, 0, false
	}
	/*gas, ok := util.CheckGasVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid gas"))
		return "", "", "", 0, 0, 0, false
	}
	nonce, ok := util.CheckNonce(req, sdk.HexToAddress(from), cdc)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid nonce"))
		return "", "", "", 0, 0, 0, false
	}
	privateKey, ok := util.CheckPrivateKey(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid private key"))
		return "", "", "", 0, 0, 0, false
	}*/
	delegatorAddrStr := from
	return from, delegatorAddrStr, privateKey, amount, gas, nonce, true
}

