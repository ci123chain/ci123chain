package rest

import (
	"encoding/hex"
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	sSDK "github.com/ci123chain/ci123chain/sdk/distribution"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func registerTxRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/distribution/tx_community_pool", rest.MiddleHandler(cliCtx, fundCommunityPoolHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/distribution/validator/withdraw_commission", rest.MiddleHandler(cliCtx, withdrawValidatorCommissionsHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/distribution/delegator/withdraw_rewards", rest.MiddleHandler(cliCtx, withdrawDelegationRewardsHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/distribution/delegator/set_withdraw_address", rest.MiddleHandler(cliCtx, setDelegatorWithdrawalAddrHandler, types.DefaultCodespace)).Methods("POST")
}

func fundCommunityPoolHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	accountAddress, ok := checkFromAddressVar(writer, req)
	if !ok {
		return
	}
	amount, ok := checkAmountVar(writer, req)
	if !ok {
		return
	}
	gas, ok := checkGasVar(writer, req)
	if !ok {
		return
	}
	nonce, err := checkNonce(writer,  req, sdk.HexToAddress(accountAddress))
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrSignTx(types.DefaultCodespace, err))
		return
	}
	privateKey, ok := checkPrivateKey(writer, req)
	if !ok {
		return
	}
	if amount.IsNegative() || amount.IsZero() {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, errors.New("invalid amount")))
		return
	}

	txByte, err := sSDK.SignFundCommunityPoolTx(accountAddress, amount, gas, nonce, privateKey)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrSignTx(types.DefaultCodespace,err))
		return
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func withdrawValidatorCommissionsHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	broadcast, err := strconv.ParseBool(req.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, req, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, err))
		return
	}
	validator := from
	msg := types.NewMsgWithdrawValidatorCommission(from, validator)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, err))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func withdrawDelegationRewardsHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	broadcast, err := strconv.ParseBool(req.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}

	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, req, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, err))
		return
	}

	validator := sdk.HexToAddress(req.FormValue("validator_address"))
	delegator := from
	msg := types.NewMsgWithdrawDelegatorReward(from, validator, delegator)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, err))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func setDelegatorWithdrawalAddrHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	broadcast, err := strconv.ParseBool(req.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, req, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, err))
		return
	}

	delegator := from
	withdraw := sdk.HexToAddress(req.FormValue("withdraw_address"))
	msg := types.NewMsgSetWithdrawAddress(from, withdraw, delegator)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, err))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}