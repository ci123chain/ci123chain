package rest

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"encoding/hex"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func SignTxRequestHandler(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		priv := request.FormValue("privateKey")
		if len(priv) < 1 {
			rest.WriteErrorResponse(writer, http.StatusNotFound, "param privateKey not found")
			return
		}

		fabric := request.FormValue("fabric")
		isFabric, err  := strconv.ParseBool(fabric)
		if err != nil {
			isFabric = false
		}
		tx, err := buildTransferTx(request, isFabric)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, errors.Wrap(err, "build transfer msg failed").Error())
			return
		}

		privPub, err := hex.DecodeString(priv)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, errors.Wrap(err, "Decode PrivateKey error").Error())
		}
		tx, err = cliCtx.SignWithTx(tx, privPub, isFabric)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, "sign tx error")
			return
		}
		txByte := tx.Bytes()
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(txByte))
	}
}

func buildTransferTx(r *http.Request, isFabric bool) (transaction.Transaction, error) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	amount := r.FormValue("amount")
	gas := r.FormValue("gas")


	froms, err := helper.ParseAddrs(from)
	if err != nil {
		return nil, err
	}
	tos, err := helper.ParseAddrs(to)
	if err != nil {
		return nil, err
	}

	if len(froms) != 1 {
		return nil, errors.New("Param from invalid")
	}
	if len(tos) != 1 {
		return nil, errors.New("Param to invalid")
	}

	gasI, err := strconv.ParseUint(gas, 10, 64)
	if err != nil {
		return nil, err
	}

	amountI, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return nil, err
	}
	nonce, err := transfer.GetNonceByAddress(froms[0])
	tx := transfer.NewTransferTx(froms[0], tos[0], gasI, nonce, types.Coin(amountI), isFabric)

	return tx, nil
}
