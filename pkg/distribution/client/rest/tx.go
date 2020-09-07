package rest

import (
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client"
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
	nonce, ok := checkNonce(writer,  req, sdk.HexToAddress(accountAddress))
	if !ok {
		return
	}
	privateKey, ok := checkPrivateKey(writer, req)
	if !ok {
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
	from, gas, nonce, priv, ok := paserArgs(writer, req)
	if !ok {
		return
	}
	validatorAddress := from

	txByte, err := sSDK.SignWithdrawValidatorCommissionTx(from, validatorAddress, gas, nonce, priv)
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func withdrawDelegationRewardsHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	from, gas, nonce, priv, ok := paserArgs(writer, req)
	if !ok {
		return
	}
	validator, ok := checkValidatorAddressVar(writer, req)
	if !ok {
		return
	}
	delegator := from
	txByte, err := sSDK.SignWithdrawDelegatorRewardTx(from, validator, delegator, gas, nonce, priv)
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

func setDelegatorWithdrawalAddrHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	from, gas, nonce, priv, ok := paserArgs(writer, req)
	if !ok {
		return
	}
	withdraw, ok := checkWithdrawAddressVar(writer, req)
	if !ok {
		return
	}
	txByte, err := sSDK.SignSetWithdrawAddressTx(from, withdraw, gas, nonce, priv)
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


func paserArgs(writer http.ResponseWriter, req *http.Request) (string, uint64, uint64, string, bool)  {
	accountAddress, ok := checkFromAddressVar(writer, req)
	if !ok {
		rest.WriteErrorRes(writer, types.ErrBadAddress(types.DefaultCodespace, errors.New(fmt.Sprintf("invalid account address: %v", accountAddress))))
		return "", 0, 0, "", false
	}
	gas, ok := checkGasVar(writer, req)
	if !ok {
		rest.WriteErrorRes(writer, types.ErrGas(types.DefaultCodespace, string(gas)))
		return "", 0, 0, "", false
	}
	nonce, ok := checkNonce(writer,  req, sdk.HexToAddress(accountAddress))
	if !ok {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, "nonce"))
		return "", 0, 0, "", false
	}
	privateKey, ok := checkPrivateKey(writer, req)
	if !ok {
		rest.WriteErrorRes(writer, types.ErrParams(types.DefaultCodespace, "privateKey"))
		return "", 0, 0, "", false
	}
	return accountAddress, gas, nonce, privateKey, true
}