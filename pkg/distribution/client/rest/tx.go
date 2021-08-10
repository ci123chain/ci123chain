package rest

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	sSDK "github.com/ci123chain/ci123chain/sdk/distribution"
	"github.com/gorilla/mux"
	"net/http"
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
		rest.WriteErrorRes(writer, err.Error())
		return
	}
	privateKey, ok := checkPrivateKey(writer, req)
	if !ok {
		return
	}
	if amount.IsNegative() || amount.IsZero() {
		rest.WriteErrorRes(writer, fmt.Sprintf("invalid amount: %v", amount))
		return
	}

	txByte, err := sSDK.SignFundCommunityPoolTx(accountAddress, amount, gas, nonce, privateKey)
	if err != nil {
		rest.WriteErrorRes(writer, err.Error())
		return
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, err.Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func withdrawValidatorCommissionsHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	validator := cliCtx.FromAddr
	msg := types.NewMsgWithdrawValidatorCommission(cliCtx.FromAddr, validator)
	if !cliCtx.Broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}
	txByte, err := types2.SignCommonTx(cliCtx.FromAddr, cliCtx.Nonce, cliCtx.Gas, []sdk.Msg{msg}, cliCtx.PrivateKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, err.Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, err.Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func withdrawDelegationRewardsHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	validator := sdk.HexToAddress(req.FormValue("validator_address"))
	delegator := cliCtx.FromAddr
	msg := types.NewMsgWithdrawDelegatorReward(cliCtx.FromAddr, validator, delegator)
	if !cliCtx.Broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(cliCtx.FromAddr, cliCtx.Nonce, cliCtx.Gas, []sdk.Msg{msg}, cliCtx.PrivateKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, err.Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, err.Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func setDelegatorWithdrawalAddrHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	delegator := cliCtx.FromAddr
	withdraw := sdk.HexToAddress(req.FormValue("withdraw_address"))
	msg := types.NewMsgSetWithdrawAddress(cliCtx.FromAddr, withdraw, delegator)
	if !cliCtx.Broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}
	txByte, err := types2.SignCommonTx(cliCtx.FromAddr, cliCtx.Nonce, cliCtx.Gas, []sdk.Msg{msg}, cliCtx.PrivateKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, err.Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, err.Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}