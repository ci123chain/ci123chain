package rest

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/gorilla/mux"
	//"github.com/ci123chain/ci123chain/pkg/order"
	"github.com/ci123chain/ci123chain/pkg/order/types"
	"net/http"
)

type OrderBook struct {
	Lists 	[]Lists 	`json:"lists"`

	Current	Current 	`json:"current"`

	Actions	[]Actions 	`json:"actions"`
}

type Lists struct {
	Name 	string 	`json:"name"`
	Height	int64	`json:"height"`
}

type Current struct {
	Index	int		`json:"index"`
	State	string	`json:"state"`
}

type Actions struct {
	Type	string	`json:"type"`
	Height	int64	`json:"height"`
	Name	string	`json:"name"`
}

/*
type ShardStateParams struct {
	Height     string    `json:"height"`
}

type QueryShardStateParams struct {
	//
	Data      ShardStateParams  `json:"data"`
}
*/

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/shared/status", QueryShardStatesRequestHandlerFn(cliCtx)).Methods("POST")
}


func QueryShardStatesRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {

		/*
		var params QueryShardStateParams
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &params)
		if readErr != nil {
			//
		}
		*/
		height := request.FormValue("height")
		prove := request.FormValue("prove")

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryState, nil, isProve)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			//rest.WriteErrorRes(writer, order.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var shardState OrderBook
		err2 := json.Unmarshal(res, &shardState)
		if err2 != nil {
			//rest.WriteErrorRes(writer, order.ErrQueryTx(types.DefaultCodespace, err2.Error()))
			return
		}
		value := shardState
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}