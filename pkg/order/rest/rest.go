package rest

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/order/types"
	sdk "github.com/ci123chain/ci123chain/sdk/shard"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {

	r.HandleFunc("/shared/add", rest.MiddleHandler(cliCtx, AddShardTxRequest, types.DefaultCodespace)).Methods("POST")
}

var cdc = app.MakeCodec()
type ShardTxBytes struct {
	From     string     `json:"from"`
	Gas      uint64     `json:"gas"`
	Nonce    uint64     `json:"nonce"`
	Type     string     `json:"type"`
	Name     string     `json:"name"`
	Height   int64     `json:"height"`
	Key      string     `json:"key"`
}

type AddShardParams struct {
	Data ShardTxBytes `json:"data"`
}

func AddShardTxRequest(cliCtx context.Context, writer http.ResponseWriter, request *http.Request) {
		key := request.FormValue("privateKey")
		gas := request.FormValue("gas")
		Gas, err := strconv.ParseInt(gas, 10, 64)
		if err != nil || Gas < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"gas error"))
			return
		}
		UserGas := uint64(Gas)
		userNonce := request.FormValue("nonce")

		from := cliCtx.GetFromAddresses()
		var nonce uint64
		if userNonce != "" {
			UserNonce, err := strconv.ParseInt(userNonce, 10, 64)
			if err != nil || UserNonce < 0 {
				rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"nonce error"))
				return
			}
			nonce = uint64(UserNonce)
		}else {
			ctx, err := client.NewClientContextFromViper(cdc)
			if err != nil {
				rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"new client context error"))
				return
			}
			nonce, _, err = ctx.GetNonceByAddress(from, false)
			if err != nil {
				rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"get nonce error"))
				return
			}
		}
		ty := request.FormValue("type")
		name := request.FormValue("name")
		height := request.FormValue("height")
		//isFabricMode := request.FormValue("isFabric")
		Height, err := strconv.ParseInt(height, 10, 64)
		if err != nil || Height < 0 {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"height error"))
			return
		}

		txByte, err := sdk.SignAddShardMsg(from,
			UserGas, nonce, ty, name, Height, key)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrCheckParams(types.DefaultCodespace,"data error"))
			return
		}

		res, err := cliCtx.BroadcastSignedData(txByte)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
}