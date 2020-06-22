package rest

import (
	"github.com/gorilla/mux"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"net/http"
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
		accountAddress := request.FormValue("accountAddress")
		height := request.FormValue("height")
		prove := request.FormValue("prove")
		checkErr := util.CheckStringLength(42, 100, accountAddress)
		if checkErr != nil {
			rest.WriteErrorRes(writer,types.ErrBadHeight(types.DefaultCodespace, checkErr))
			return
		}

		if height == "" {
			height = "now"
		}

		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/rewards/" + accountAddress + "/" + height, nil, isProve)
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
		value := &RewardsData{Rewards:rewards}
		if height == "now" {
			height = ""
		}
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}