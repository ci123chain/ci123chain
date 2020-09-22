package rest

import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client"
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

var cdc = app.MakeCodec()

func CreateValidatorRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	delegatorAddr := from
	validatorAddr := from
	amount, err := strconv.ParseUint(request.FormValue("amount"), 10, 64)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}

	msd, r, mr, mcr, moniker, identity, website, securityContact, details,
	publicKey,err := parseOtherArgs(request)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	/*by, err := hex.DecodeString(publicKey)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	var public crypto.PubKey
	err = cdc.UnmarshalJSON(by, &public)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}*/

	coin := sdk.NewUInt64Coin(amount)
	MSD, R, MR, MXR := sSdk.CreateParseArgs(msd, r, mr, mcr)
	msg := staking.NewCreateValidatorMsg(from, coin, MSD, validatorAddr,
		delegatorAddr, R, MR, MXR, moniker, identity, website, securityContact, details, publicKey)

	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
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
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	delegatorAddr := from
	validatorAddr := sdk.HexToAddress(request.FormValue("validator_address"))
	amount, err := strconv.ParseUint(request.FormValue("amount"), 10, 64)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from, validatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}

	coin := sdk.NewUInt64Coin(amount)
	msg := staking.NewDelegateMsg(from, delegatorAddr, validatorAddr, coin)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
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
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	delegatorAddr := from
	validatorSrcAddr := sdk.HexToAddress(request.FormValue("validator_src_address"))
	validatorDstAddr := sdk.HexToAddress(request.FormValue("validator_dst_address"))
	amount, err := strconv.ParseUint(request.FormValue("amount"), 10, 64)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from, validatorSrcAddr, validatorDstAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}

	coin := sdk.NewUInt64Coin(amount)
	msg := staking.NewRedelegateMsg(from, delegatorAddr, validatorSrcAddr, validatorDstAddr, coin)

	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
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
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	delegatorAddr := from
	validatorAddr := sdk.HexToAddress(request.FormValue("validator_address"))
	amount, err := strconv.ParseUint(request.FormValue("amount"), 10, 64)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from, validatorAddr)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}

	coin := sdk.NewUInt64Coin(amount)
	msg := staking.NewUndelegateMsg(from, delegatorAddr, validatorAddr, coin)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
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
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
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
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, "new_rate error"))
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
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, "min_self_delegation error"))
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

	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
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
	ms := req.FormValue("min_self_delegation")
	if ms != "" {
		minSelf, err = util.CheckInt64(ms)
		if err != nil || minSelf <= 0 {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid min_self_delegation"))
			return 0, 0, false
		}
	}else {
		minSelf = -1
	}
	nr := req.FormValue("new_rate")
	if nr != "" {
		newRate, err = util.CheckInt64(nr)
		if err != nil || newRate > 100 || newRate <= 0 {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid new_rate"))
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

/*func parseBasicArgs(w http.ResponseWriter, req *http.Request) (string, uint64, uint64, string, bool) {

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
}*/

/*func checkNonce(r *http.Request, from sdk.AccAddress, cdc *codec.Codec) (uint64, bool) {
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
}*/


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



/*
func parseBaseArgs(w http.ResponseWriter, req *http.Request) (string, string, string, int64, uint64, uint64, bool) {

	//为了确保 createValidator交易是from账户发出的，validator的地址就直接是这个from;
	//只能是 设置 提取奖金的地址；那么就要保证 from validator delegator是一致的；但是 又没有验证；因为它本身就是使用相同的值；
	//要保证delegator跟from一致；
	from, gas, nonce, privateKey, ok := parseBasicArgs(w, req)
	if !ok {
		return "", "", "", 0, 0, 0, false
	}
	amount, ok := util.CheckAmountVar(req)
	if !ok {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "invalid amount"))
		return "", "", "", 0, 0, 0, false
	}
	delegatorAddrStr := from
	return from, delegatorAddrStr, privateKey, amount, gas, nonce, true
}
*/

