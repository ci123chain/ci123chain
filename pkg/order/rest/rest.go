package rest

import (
	"encoding/hex"
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/order/types"
	"net/http"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {

	r.HandleFunc("/tx/addShard", AddShardTxRequest(cliCtx)).Methods("POST")
}

type ShardTxBytes struct {
	Data string `json:"data"`
}

func AddShardTxRequest(cliCtx context.Context) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {

		/*var shardTxBytes ShardTxBytes
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &shardTxBytes)
		if readErr != nil {
			//
		}*/
		data := request.FormValue("data")
		txByte, err := hex.DecodeString(data)
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