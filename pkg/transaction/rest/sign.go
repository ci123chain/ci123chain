package rest

import (
	"CI123Chain/pkg/abci/types/rest"
	"CI123Chain/pkg/client/context"
	"CI123Chain/pkg/client/helper"
	"CI123Chain/pkg/transaction"
	"encoding/hex"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func SignTxRequestHandler(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		txType := request.FormValue("code")
		code, err := strconv.Atoi(txType)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusNotFound, "param code not found")
			return
		}
		priv := request.FormValue("privateKey")
		if len(priv) < 1 {
			rest.WriteErrorResponse(writer, http.StatusNotFound, "param privateKey not found")
			return
		}
		if uint8(code) == transaction.TRANSFER {
			tx, err := buildTransferTx(request)
			if err != nil {
				rest.WriteErrorResponse(writer, http.StatusNotFound, errors.Wrap(err, "build transfer msg failed").Error())
				return
			}

			tx, err = cliCtx.SignTx2(tx, priv)
			if err != nil {
				rest.WriteErrorResponse(writer, http.StatusNotFound, "sign tx error")
				return
			}
			txByte := tx.Bytes()
			//fmt.Println(txByte)
			rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(txByte))
		}
	}
}

func buildTransferTx(r *http.Request) (transaction.Transaction, error) {
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
	gasI, err := strconv.ParseUint(gas, 10, 64)
	if err != nil {
		return nil, err
	}

	amountI, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return nil, err
	}
	nonce, err := transaction.GetNonceByAddress(froms[0])
	tx := transaction.NewTransferTx(froms[0], tos[0], gasI, nonce, amountI)

	return tx, nil
}
