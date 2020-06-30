package rest

import (
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
	r.HandleFunc("/distribution/tx_community_pool", fundCommunityPoolHandlerFn(cliCtx)).Methods("POST")
}


func fundCommunityPoolHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
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

		res, err := cliCtx.BroadcastSignedData(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
	}
}