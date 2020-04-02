package rest

import (
	"encoding/hex"
	"github.com/pkg/errors"
	//sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/transfer/types"
	tSDK "github.com/tanhuiya/ci123chain/sdk/transfer"
	"net/http"
	"strconv"
)

var cdc = app.MakeCodec()

type Tx struct {
	SignedTx	string `json:"signedtx"`
}


func SignTxRequestHandler(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		priv := request.FormValue("privateKey")
		if len(priv) < 1 {
			rest.WriteErrorRes(writer, transaction.ErrBadPrivkey(types.DefaultCodespace, errors.New("param privateKey not found")) )
			return
		}
		fabric := request.FormValue("fabric")
		isFabric, err  := strconv.ParseBool(fabric)
		if err != nil {
			isFabric = false
		}
		txByte, err := buildTransferTx(request, isFabric, priv)
		if err != nil {
			rest.WriteErrorRes(writer, transaction.ErrSignature(types.DefaultCodespace, errors.New("sign with tx error")))
			return
		}
		resp := &Tx{SignedTx:hex.EncodeToString(txByte)}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

func buildTransferTx(r *http.Request, isFabric bool, priv string) ([]byte, error) {

	fabric := r.FormValue("fabric")
	isFabric, err  := strconv.ParseBool(fabric)
	if err != nil {
		isFabric = false
	}
	from := r.FormValue("from")
	to := r.FormValue("to")
	amount := r.FormValue("amount")
	gas := r.FormValue("gas")
	userNonce := r.FormValue("nonce")


	froms, err := helper.ParseAddrs(from)
	if err != nil {
		return nil, client.ErrParseAddr(types.DefaultCodespace, err)
	}
	if len(froms) != 1 {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "from error")
	}
	/*
		tos, err := helper.ParseAddrs(to)
		if err != nil {
			return nil, client.ErrParseAddr(types.DefaultCodespace, err)
		}
		if len(tos) != 1 {
			return nil, types.ErrCheckParams(types.DefaultCodespace, "to error")
		}
	*/
	gasI, err := strconv.ParseUint(gas, 10, 64)
	if err != nil {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "gas error")
	}

	amountI, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "amount error")
	}
	var nonce uint64
	if userNonce != "" {
		UserNonce, err := strconv.ParseInt(userNonce, 10, 64)
		if err != nil || UserNonce < 0 {
			return nil, types.ErrCheckParams(types.DefaultCodespace, "nonce error")
		}
		nonce = uint64(UserNonce)
	}else {
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return nil, client.ErrNewClientCtx(types.DefaultCodespace, err)
		}
		nonce, err = ctx.GetNonceByAddress(froms[0])
		if err != nil {
			return nil, types.ErrCheckParams(types.DefaultCodespace, "nonce error")
		}
	}
	tx, err := tSDK.SignTransferMsg(from, to, amountI, gasI, nonce, priv, isFabric)
	if err != nil {
		return nil, err
	}

	return tx, nil
}