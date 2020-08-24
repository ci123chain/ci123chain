package rest

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transfer/rest/utils"
	"github.com/ci123chain/ci123chain/pkg/util"
	"net/http"
)

type TxRequestParams struct {
	Hash    string    `json:"hash"`
	Height  string    `json:"height"`
}

func QueryTxRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		hashHexStr := request.FormValue("hash")
		checkErr := util.CheckStringLength(1, 100, hashHexStr)
		if checkErr != nil {
			rest.WriteErrorRes(writer, types.ErrQueryTx(types.DefaultCodespace, checkErr.Error()))
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		resp, err := utils.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if resp.Empty() {
			rest.WriteErrorRes(writer, types.ErrQueryTx(types.DefaultCodespace,fmt.Sprintf("no transfer found with hash %s", hashHexStr)))
		}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}