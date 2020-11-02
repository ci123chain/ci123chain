package rest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/infrastructure"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func RegisterRestTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/infrastructure/store_content", rest.MiddleHandler(cliCtx, StoreContentRequest, types.DefaultCodespace)).Methods("POST")
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
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}

	//verify account exists
	err = checkAccountExist(cliCtx, from)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	key_str := request.FormValue("key")
	content_str := request.FormValue("content")
	if key_str == ""{
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, "key"))
		return
	}
	if content_str == ""{
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace, "content"))
		return
	}
	//key, err := hex.DecodeString(key_str)
	//content, err := hex.DecodeString(content_str)
	value, err := json.Marshal(types.NewStoredContent(key_str, content_str))
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrMarshalFailed(types.DefaultCodespace, err))
		return
	}
	msg := infrastructure.NewStoreContentMsg(from, key_str, value)

	if !broadcast {
		rest.PostProcessResponseBare(writer, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
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