package rest

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/slashing/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func registerTxHandlers(clientCtx context.Context, r *mux.Router) {
	r.HandleFunc("/slashing/validators/unjail", rest.MiddleHandler(clientCtx, NewUnjailRequestHandlerFn, types.DefaultCodespace)).Methods("POST")
}
//
//// Unjail TX body
//type UnjailReq struct {
//	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
//}
//

var cdc = types2.GetCodec()

// NewUnjailRequestHandlerFn returns an HTTP REST handler for creating a MsgUnjail
// transaction.
func NewUnjailRequestHandlerFn(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}

	privKey, fromAddr, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)

	msg := types.NewMsgUnjail(fromAddr)
	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(fromAddr, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}
