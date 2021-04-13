package rest

import (
	"encoding/hex"
	abcitypes "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/order/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/shared/add", rest.MiddleHandler(cliCtx, AddShardTxRequest, sdkerrors.RootCodespace)).Methods("POST")
}

var cdc = types2.GetCodec()

func AddShardTxRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	broatcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broatcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broatcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	Type := request.FormValue("types")
	name := request.FormValue("name")
	height := request.FormValue("height")
	//isFabricMode := request.FormValue("isFabric")
	Height, err := strconv.ParseInt(height, 10, 64)
	if err != nil || Height < 0 {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid height").Error())
		return
	}

	msg := types.NewMsgUpgrade(from, Type, name, Height)
	if !broatcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []abcitypes.Msg{msg}, privKey, cdc)
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