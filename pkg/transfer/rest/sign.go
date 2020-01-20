package rest

import (
	"encoding/hex"
	"github.com/pkg/errors"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"github.com/tanhuiya/ci123chain/pkg/transfer/types"
	"net/http"
	"strconv"
)

var cdc = app.MakeCodec()

type Tx struct {
	SignedTx	string `json:"signedtx"`
}

/*type TxAccountParams struct {
	From       string    `json:"from"`
	To         string    `json:"to"`
	Gas        string    `json:"gas"`
	Amount     string    `json:"amount"`
	Key        string    `json:"key"`
	Fabric     string     `json:"fabric"`
}

type TxParams struct {
	Data TxAccountParams `json:"data"`
}*/


func SignTxRequestHandler(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		/*
		var params TxParams
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &params)
		if readErr != nil {
			//
		}
		*/
		//priv := params.Data.Key
		priv := request.FormValue("privateKey")
		if len(priv) < 1 {
			rest.WriteErrorRes(writer, transaction.ErrBadPrivkey(types.DefaultCodespace, errors.New("param privateKey not found")) )
			return
		}
		/*
		from := params.Data.From
		to := params.Data.To
		gas := params.Data.Gas
		amount := params.Data.Amount

		fabric := params.Data.Fabric
		*/
		fabric := request.FormValue("fabric")
		isFabric, err  := strconv.ParseBool(fabric)
		if err != nil {
			isFabric = false
		}
		tx, err := buildTransferTx(request, isFabric)
		if err != nil {
			rest.WriteErrorRes(writer, err.(sdk.Error))
			return
		}

		privPub, err := hex.DecodeString(priv)
		if err != nil {
			rest.WriteErrorRes(writer, transaction.ErrBadPrivkey(types.DefaultCodespace, err))
		}
		tx, err = cliCtx.SignWithTx(tx, privPub, isFabric)
		if err != nil {
			rest.WriteErrorRes(writer, transaction.ErrSignature(types.DefaultCodespace, errors.New("sign with tx error")))
			return
		}
		txByte := tx.Bytes()
		resp := &Tx{SignedTx:hex.EncodeToString(txByte)}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

func buildTransferTx(r *http.Request, isFabric bool) (transaction.Transaction, error) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	amount := r.FormValue("amount")
	gas := r.FormValue("gas")


	froms, err := helper.ParseAddrs(from)
	if err != nil {
		return nil, client.ErrParseAddr(types.DefaultCodespace, err)
	}
	tos, err := helper.ParseAddrs(to)
	if err != nil {
		return nil, client.ErrParseAddr(types.DefaultCodespace, err)
	}

	if len(froms) != 1 {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "from error")
	}
	if len(tos) != 1 {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "to error")
	}

	gasI, err := strconv.ParseUint(gas, 10, 64)
	if err != nil {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "gas error")
	}

	amountI, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return nil, types.ErrCheckParams(types.DefaultCodespace, "amount error")
	}
	ctx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return nil, client.ErrNewClientCtx(types.DefaultCodespace, err)
	}
	nonce, err := ctx.GetNonceByAddress(froms[0])
	tx := transfer.NewTransferTx(froms[0], tos[0], gasI, nonce, sdk.NewUInt64Coin(amountI), isFabric)

	return tx, nil
}
