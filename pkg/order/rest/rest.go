package rest

import (
	"encoding/hex"
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/order/types"
	order "github.com/tanhuiya/ci123chain/pkg/order/types"
	"net/http"
	"strconv"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {

	r.HandleFunc("/tx/addShard", AddShardTxRequest(cliCtx)).Methods("POST")
}

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

func AddShardTxRequest(cliCtx context.Context) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {

		/*
		var shardTxBytes AddShardParams
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &shardTxBytes)
		if readErr != nil {
			//
		}
		*/
		data := request.FormValue("data")
		from := request.FormValue("from")
		gas := request.FormValue("gas")
		Gas, err := strconv.ParseInt(gas, 10, 64)
		UserGas := uint64(Gas)
		nonce := request.FormValue("nonce")
		Nonce, err := strconv.ParseInt(nonce, 10, 64)
		UserNonce := uint64(Nonce)
		ty := request.FormValue("type")
		name := request.FormValue("name")
		height := request.FormValue("height")
		Height, err := strconv.ParseInt(height, 10, 64)
		//UserHeight := uint64(Height)
		privByte, err := hex.DecodeString(data)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
			return
		}

		txByte, err := order.SignUpgradeTx(from,
			UserGas, UserNonce, ty, name, Height, privByte)
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
}