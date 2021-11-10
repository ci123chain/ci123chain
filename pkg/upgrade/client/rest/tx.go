package rest

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/upgrade/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// RegisterTxRoutes - Central function to define routes that get registered by the main application
func RegisterTxRoutes(clientCtx context.Context, r *mux.Router) {
	r.HandleFunc("/upgrade/proposal", rest.MiddleHandler(clientCtx, NewUpgradeProposalHandlerFn, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/upgrade/cancel_proposal", rest.MiddleHandler(clientCtx, NewUpgradeCancelProposalHandlerFn, types.DefaultCodespace)).Methods("POST")
}

func NewUpgradeProposalHandlerFn(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	privKey, fromAddr, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cliCtx.Cdc, true)

	name := request.FormValue("name")
	desc := request.FormValue("desc")
	heightStr := request.FormValue("height")
	height, err := strconv.ParseInt(heightStr, 10, 64)

	plan := types.Plan{
		Name:   name,
		Height: height,
		Info:   desc,
	}

	msg := types.NewSoftwareUpgradeProposal(fromAddr, name, desc, plan)

	txByte, err := types2.SignCommonTx(fromAddr, nonce, gas, []sdk.Msg{msg}, privKey, cliCtx.Cdc)
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

func NewUpgradeCancelProposalHandlerFn(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	privKey, fromAddr, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cliCtx.Cdc, true)
	msg := types.NewCancelSoftwareUpgradeProposal(fromAddr)
	txByte, err := types2.SignCommonTx(fromAddr, nonce, gas, []sdk.Msg{msg}, privKey, cliCtx.Cdc)
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