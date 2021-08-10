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
	Type := request.FormValue("type")
	name := request.FormValue("name")
	height := request.FormValue("height")
	//isFabricMode := request.FormValue("isFabric")
	Height, err := strconv.ParseInt(height, 10, 64)
	if err != nil || Height < 0 {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid height").Error())
		return
	}

	msg := types.NewMsgUpgrade(cliCtx.FromAddr, Type, name, Height)
	if !cliCtx.Broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(cliCtx.FromAddr, cliCtx.Nonce, cliCtx.Gas, []abcitypes.Msg{msg}, cliCtx.PrivateKey, cdc)
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