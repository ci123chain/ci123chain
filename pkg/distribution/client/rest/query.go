package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/distribution/types"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"io/ioutil"
	"net/http"
	"strconv"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/rewards", QueryValidatorRewardsRequestHandlerFn(cliCtx)).Methods("POST")
}

type RewardsData struct {
	Rewards 	uint64 `json:"rewards"`
}

type RewardsParams struct {
	Address string `json:"address"`
	Height  string     `json:"height"`
}

type QueryRewardsParams struct {
	Data RewardsParams `json:"data"`
}

func QueryValidatorRewardsRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		//accountAddress := vars["accountAddress"]
		//height := vars["height"]

		var params QueryRewardsParams
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &params)
		if readErr != nil {
			//
		}

		if params.Data.Height == "now" {
		}else {
			_, Err := strconv.ParseInt(params.Data.Height, 10 , 64)
			if Err != nil {
				rest.WriteErrorRes(writer,types.ErrBadHeight(types.DefaultCodespace, Err))
				return
			}
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/rewards/" + params.Data.Address + "/" + params.Data.Height, nil)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var rewards uint64
		err2 := cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(res, &rewards)
		if err2 != nil {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, err2.Error()))
			return
		}
		resp := &RewardsData{Rewards:rewards}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}