package rest

import (
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transfer/rest/utils"
	"github.com/ci123chain/ci123chain/pkg/util"
	"net/http"
	"strings"
)

type TxRequestParams struct {
	Hash    string    `json:"hash"`
	Height  string    `json:"height"`
}

func QueryTxRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		hashHexStr := request.FormValue("hash")
		if strings.HasPrefix(hashHexStr, "0x") {
			hashHexStr= hashHexStr[2:]
		}
		checkErr := util.CheckStringLength(1, 100, hashHexStr)
		if checkErr != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("invalid hash")).Error())
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, "get clictx failed")
			return
		}

		resp, err := utils.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if resp.Empty() {
			rest.WriteErrorRes(writer,  fmt.Sprintf("no tx found with hash %s", hashHexStr))
		}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

func QueryTxsWithHeight(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		heightStr := request.FormValue("heights")
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, "get clictx failed")
			return
		}
		var heights sdk.Heights
		Err := json.Unmarshal([]byte(heightStr), &heights)
		if Err != nil {
			rest.WriteErrorRes(writer, fmt.Sprintf("invalid heights %v", heightStr))
			return
		}

		resp, err := utils.QueryTxsWithHeight(cliCtx, heights)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}else {
			rest.PostProcessResponseBare(writer, cliCtx, resp)
		}
	}
}