package rest

import (
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/order/types"
	order "github.com/tanhuiya/ci123chain/pkg/order/types"
	"io/ioutil"
	"net/http"
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

		var shardTxBytes AddShardParams
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &shardTxBytes)
		if readErr != nil {
			//
		}
		//data := request.FormValue("data")
		privByte, err := hex.DecodeString(shardTxBytes.Data.Key)
		if err != nil {
			rest.WriteErrorRes(writer, client.ErrBroadcast(types.DefaultCodespace, err))
			return
		}

		txByte, err := order.SignUpgradeTx(shardTxBytes.Data.From,
			shardTxBytes.Data.Gas, shardTxBytes.Data.Nonce, shardTxBytes.Data.Type, shardTxBytes.Data.Name, shardTxBytes.Data.Height, privByte)
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