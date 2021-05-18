package rest

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

var (
	DefaultVersion uint64 = 0
)


func RegisterTxRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/ibc/transfer", rest.MiddleHandler(cliCtx, ibcTransferHandler, "")).Methods("POST")
}


func ibcTransferHandler(cliCtx context.Context, writer http.ResponseWriter, req *http.Request) {
	sourcePort := req.FormValue("source_port")
	sourceChannel := req.FormValue("source_channel")
	token := req.FormValue("token")
	coin, err := sdk.ParseCoinNormalized(token)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	receiver := req.FormValue("receiver")
	timeoutHeight, err := strconv.ParseInt(req.FormValue("timeout_height"), 10, 64)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	timeoutTimestamp, err := strconv.ParseInt(req.FormValue("timeout_timestamp_offset"), 10, 64)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	if timeoutTimestamp != 0 {
		timeoutTimestamp = time.Now().UTC().Add(time.Duration(timeoutTimestamp) * time.Second).UnixNano()
	}

	privateKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, req, cliCtx.Cdc, true)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	if coin.IsNegative() || coin.IsZero() {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid coin").Error())
		return
	}

	msg := types.NewMsgTransfer(
		sourcePort,
		sourceChannel,
		coin,
		from.String(), // 'MustGetAddress' must be called directly before calling 'NewMsg...'
		receiver,
		clienttypes.NewHeight(DefaultVersion, uint64(timeoutHeight)),
		uint64(timeoutTimestamp),
	)
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid coin").Error())
		return
	}
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privateKey, cliCtx.Cdc)

	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "broadcast tx failed").Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}
