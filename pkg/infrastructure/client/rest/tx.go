package rest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/infrastructure"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func RegisterRestTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/infrastructure/store_content", rest.MiddleHandler(cliCtx, StoreContentRequest, sdkerrors.RootCodespace)).Methods("POST")
}

var cdc = types2.MakeCodec()


func StoreContentRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
	//
	broadcast, err := strconv.ParseBool(request.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, request, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error()).Error())
		return
	}
	key_str := request.FormValue("key")
	content_str := request.FormValue("content")
	if key_str == ""{
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrParams, "key cant not be empty").Error())
		return
	}
	if content_str == ""{
		rest.WriteErrorRes(writer,sdkerrors.Wrap(sdkerrors.ErrParams, "conten can not be empty").Error())
		return
	}
	//key, err := hex.DecodeString(key_str)
	//content, err := hex.DecodeString(content_str)
	value, err := json.Marshal(types.NewStoredContent(key_str, content_str))
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error()).Error())
		return
	}
	msg := infrastructure.NewStoreContentMsg(from, key_str, value)

	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "sign tx failed").Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "broadcast tx failed").Error())
		return
	}
	rest.PostProcessResponseBare(writer, cliCtx, res)
}

func checkAccountExist(ctx context.Context, address... sdk.AccAddress) error {
	for i := 0; i < len(address); i++ {
		_,_,  err := ctx.GetNonceByAddress(address[i], false)
		if err != nil {
			return errors.New(fmt.Sprintf("account of %s does not exist", address[i]))
		}
	}
	return nil
}