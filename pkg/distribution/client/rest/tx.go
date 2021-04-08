package rest

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	privateKey, ok := checkPrivateKey(writer, req)
	if !ok {
		return
	}
	if amount.IsNegative() || amount.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid amount").Error())
		return
	}

	txByte, err := sSDK.SignFundCommunityPoolTx(accountAddress, amount, gas, nonce, privateKey)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "sign tx failed").Error())
		return
	}

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "broadcast tx failed").Error())
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "sign tx failed").Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "broadcast tx failed").Error())
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "sign tx failed").Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "broadcast tx failed").Error())
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
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
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "sign tx failed").Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "broadcast tx failed").Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}