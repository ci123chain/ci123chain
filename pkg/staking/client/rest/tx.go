package rest

import (
	"errors"
	"fmt"
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

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/staking/validator/create", CreateValidatorRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/delegate", DelegateTX(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/redelegate", RedelegateTX(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/undelegate", UndelegateTX(cliCtx)).Methods("POST")
}

var cdc = app.MakeCodec()


func CreateValidatorRequest(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

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
		validatorAddr := request.FormValue("validatorAddress")
		//verify account exists
		err = checkAccountExist(cliCtx, from, validatorAddr)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		}


		txByte, err := sSdk.SignCreateValidatorMSg(from, amount, gas, nonce, key, msd, validatorAddr,
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

		res, err := cliCtx.BroadcastSignedData(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
	}
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
	from, ok := util.CheckFromAddressVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid from address"))
		return "", "", "", 0, 0, 0, false
	}
	amount, ok := util.CheckAmountVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid amount"))
		return "", "", "", 0, 0, 0, false
	}
	gas, ok := util.CheckGasVar(req)
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
	}
	delegatorAddrStr := from
	return from, delegatorAddrStr, privateKey, amount, gas, nonce, true
}

func checkAccountExist(ctx context.Context, address... string) error {
	for i := 0; i < len(address); i++ {
		_, err := ctx.GetNonceByAddress(sdk.HexToAddress(address[i]))
		if err != nil {
			return errors.New(fmt.Sprintf("account of %s does not exist", address[i]))
		}
	}
	/*for _, addr := range address {
		_, err := ctx.GetNonceByAddress(sdk.HexToAddress(addr))
		if err != nil {
			return errors.New(fmt.Sprintf("account of %s does not exist", addr))
		}
	}*/
	return nil
}