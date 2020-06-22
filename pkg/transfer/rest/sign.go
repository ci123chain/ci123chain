package rest

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/ci123chain/ci123chain/pkg/util"

	//sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	tSDK "github.com/ci123chain/ci123chain/sdk/transfer"
	"net/http"
)

var cdc = app.MakeCodec()

type Tx struct {
	SignedTx	string `json:"signedtx"`
}

func SignTxRequestHandler(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		priv := request.FormValue("privateKey")
		err := util.CheckStringLength(1, 100, priv)
		if err != nil {
			rest.WriteErrorRes(writer, transaction.ErrBadPrivkey(types.DefaultCodespace, errors.New("param privateKey not found")) )
			return
		}
		fabric := request.FormValue("fabric")
		isFabric, err := util.CheckBool(fabric)
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

	from := r.FormValue("from")
	to := r.FormValue("to")
	amount := r.FormValue("amount")
	gas := r.FormValue("gas")
	userNonce := r.FormValue("nonce")


	froms, err := helper.ParseAddrs(from)
	if err != nil {
		return nil, client.ErrParseAddr(types.DefaultCodespace, err)
	}
	err = util.CheckStringLength(42, 100, from)
	if err != nil {
		return nil, err
	}
	err = util.CheckStringLength(42, 100, to)
	if err != nil {
		return nil, err
	}

	gasI, err := util.CheckUint64(gas)
	if err != nil {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "gas error")
	}

	amountI, err := util.CheckUint64(amount)
	if err != nil {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "amount error")
	}
	var nonce uint64
	if userNonce != "" {
		UserNonce, err := util.CheckUint64(userNonce)
		if err != nil || UserNonce < 0 {
			return nil, types.ErrCheckParams(types.DefaultCodespace, "nonce error")
		}
		nonce = UserNonce
	}else {
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return nil, client.ErrNewClientCtx(types.DefaultCodespace, err)
		}
		nonce, _, err = ctx.GetNonceByAddress(froms[0],false)
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