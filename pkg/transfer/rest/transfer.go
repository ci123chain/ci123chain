package rest

import (
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/transfer/types"
	"io/ioutil"
	"net/http"
)

func SendRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		var params Params
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &params)
		if readErr != nil {
			//
		}

		priv := params.Data
		//priv := request.FormValue("privateKey")
		if len(priv) < 1 {
			rest.WriteErrorRes(writer, transaction.ErrBadPrivkey(types.DefaultCodespace, errors.New("param privateKey not found")) )
			return
		}

		tx, err := buildTransferTx(request, false)
		if err != nil {
			rest.WriteErrorRes(writer, err.(sdk.Error))
			return
		}

		privPub, err := hex.DecodeString(priv)
		if err != nil {
			rest.WriteErrorRes(writer, transaction.ErrBadPrivkey(types.DefaultCodespace, err))
		}
		tx, err = cliCtx.SignWithTx(tx, privPub, false)
		if err != nil {
			rest.WriteErrorRes(writer, transaction.ErrSignature(types.DefaultCodespace, errors.New("sign with tx error")))
			return
		}
		txByte := tx.Bytes()

		res, err := cliCtx.BroadcastSignedData(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
	}
}

