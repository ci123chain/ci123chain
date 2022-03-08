package rest

import (
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterQueryRoutes(cliCtx context.Context, r *mux.Router) {
	// Get all validators
	r.HandleFunc("/preStaking/stakingRecord", QueryStakingRecord(cliCtx), ).Methods("POST")
	r.HandleFunc("/preStaking/getDao", QueryPreStakingDao(cliCtx))
}

func QueryPreStakingDao(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		//
		res, _, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.PreStakingTokenQuery, nil, false)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, fmt.Sprintf("unexpected res: %v", res))
			return
		}
		var result string
		err = json.Unmarshal(res, &result)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		value := result
		resp := rest.BuildQueryRes("", false, value, nil)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}


func QueryStakingRecord(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		delegatorAddress := req.FormValue("delegator_address")
		prove := req.FormValue("prove")
		height := req.FormValue("height")
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}
		params := types.QueryStakingRecord{
			DelegatorAddr:delegatorAddr,
		}
		bz, err := cliCtx.Cdc.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorRes(writer, fmt.Sprintf("cdc marshal failed: %v", err.Error()))
			return
		}
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.StakingRecordQuery, bz, isProve)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, fmt.Sprintf("unexpected res: %v", res))
			return
		}

		var result []types.StakingVault
		err = types.PreStakingCodec.UnmarshalBinaryBare(res, &result)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		value := result
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}