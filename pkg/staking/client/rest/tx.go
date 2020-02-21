package rest

import (
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
	sSdk "github.com/tanhuiya/ci123chain/sdk/staking"
	"net/http"
	"strconv"
)

func RegisterRestTxRoutes(cliCtx context.Context, r *mux.Router)  {

	r.HandleFunc("/staking/delegate", DelegateTX(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/redelegate", RedelegateTX(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/undelegate", UndelegateTX(cliCtx)).Methods("POST")
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
		UserGas := uint64(Gas)
		nonce := request.FormValue("nonce")
		Nonce, err := strconv.ParseInt(nonce, 10, 64)
		if err != nil || Nonce < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"nonce error"))
			return
		}
		UserNonce := uint64(Nonce)
		amount := request.FormValue("amount")
		amt, err := strconv.ParseInt(amount, 10, 64)
		if err != nil || Nonce < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"amount of coin error"))
			return
		}
		cAmt := uint64(amt)
		validatorAddr := request.FormValue("validatorAddr")
		delegatorAddr := request.FormValue("delegatorAddr")

		txByte, err := sSdk.SignDelegateMsg(from,cAmt, UserGas, UserNonce, key, validatorAddr, delegatorAddr)
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
		nonce := request.FormValue("nonce")
		Nonce, err := strconv.ParseInt(nonce, 10, 64)
		if err != nil || Nonce < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"nonce error"))
			return
		}
		UserNonce := uint64(Nonce)
		amount := request.FormValue("amount")
		amt, err := strconv.ParseInt(amount, 10, 64)
		if err != nil || Nonce < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"amount of coin error"))
			return
		}
		cAmt := uint64(amt)
		validatorSrcAddr := request.FormValue("validatorSrcAddr")
		validatorDstAddr := request.FormValue("validatorDstAddr")
		delegatorAddr := request.FormValue("delegatorAddr")

		txByte, err := sSdk.SignRedelegateMsg(from,cAmt, UserGas, UserNonce, key, validatorSrcAddr, validatorDstAddr, delegatorAddr)
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
		nonce := request.FormValue("nonce")
		Nonce, err := strconv.ParseInt(nonce, 10, 64)
		if err != nil || Nonce < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"nonce error"))
			return
		}
		UserNonce := uint64(Nonce)
		amount := request.FormValue("amount")
		amt, err := strconv.ParseInt(amount, 10, 64)
		if err != nil || Nonce < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"amount of coin error"))
			return
		}
		cAmt := uint64(amt)
		validatorAddr := request.FormValue("validatorAddr")
		delegatorAddr := request.FormValue("delegatorAddr")

		txByte, err := sSdk.SignUndelegateMsg(from,cAmt, UserGas, UserNonce, key, validatorAddr, delegatorAddr)
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