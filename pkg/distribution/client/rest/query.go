package rest

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/distribution/client/common"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterQueryRoutes(cliCtx context.Context, r *mux.Router)  {
	//r.HandleFunc("/rewards", QueryValidatorRewardsRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/distribution/validator/outstanding_rewards", QueryValidatorOutstandingRewardsHandleFn(cliCtx)).Methods("POST")
	r.HandleFunc("/distribution/query_community_pool", QueryCommunityPoolHandleFn(cliCtx)).Methods("POST")
	r.HandleFunc("/distribution/delegator/withdraw_address", QueryWithDrawAddress(cliCtx)).Methods("POST")
	r.HandleFunc("/distribution/validator/rewards", validatorInfoHandleFn(cliCtx)).Methods("POST")
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

/*func QueryValidatorRewardsRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		accountAddress := vars["accountAddress"]
		height := vars["height"]
		checkErr := util.CheckStringLength(42, 100, accountAddress)
		if checkErr != nil {
			rest.WriteErrorRes(writer,types.ErrBadAddress(types.DefaultCodespace, checkErr))
			return
		}

		if height == "" {
			height = "now"
		}else {
			_, Err := util.CheckInt64(height)
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

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/rewards/" + accountAddress + "/" + height, nil)
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
}*/

func QueryValidatorOutstandingRewardsHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		validatorAddress, ok := checkValidatorAddressVar(writer, request)
		if !ok {
			return
		}
		b := cliCtx.Cdc.MustMarshalJSON(types.NewQueryValidatorOutstandingRewardsParams(validatorAddress))
		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryValidatorOutstandingRewards, b)
		if err != nil {
			rest.WriteErrorRes(writer, err)
		}
		var rewards types.ValidatorOutstandingRewards
		cliCtx.Cdc.MustUnmarshalJSON(res, &rewards)
		rest.PostProcessResponseBare(writer, cliCtx, rewards)
	}
}

func QueryCommunityPoolHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryCommunityPool, nil)
		if err != nil {
			rest.WriteErrorRes(w, err)
		}
		var result sdk.DecCoin
		cliCtx.Cdc.MustUnmarshalJSON(res, &result)
		rest.PostProcessResponseBare(w, cliCtx, result)
	}
}

func QueryWithDrawAddress(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		delegatorAddr, ok := checkDelegatorAddressVar(w, r)
		if !ok {
			return
		}
		b := cliCtx.Cdc.MustMarshalJSON(types.NewQueryDelegatorWithdrawAddrParams(delegatorAddr))
		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryWithdrawAddress, b)
		if err != nil {
			rest.WriteErrorRes(w, err)
		}
		var result sdk.AccAddress
		cliCtx.Cdc.MustUnmarshalJSON(res, &result)
		rest.PostProcessResponseBare(w, cliCtx, result)
	}
}

// ValidatorDistInfo defines the properties of
// validator distribution information response.
type ValidatorDistInfo struct {
	OperatorAddress     sdk.AccAddress                       `json:"operator_address" yaml:"operator_address"`
	SelfBondRewards     sdk.DecCoin                         `json:"self_bond_rewards" yaml:"self_bond_rewards"`
	ValidatorCommission types.ValidatorAccumulatedCommission `json:"val_commission" yaml:"val_commission"`
}

// NewValidatorDistInfo creates a new instance of ValidatorDistInfo.
func NewValidatorDistInfo(operatorAddr sdk.AccAddress, rewards sdk.DecCoin,
	commission types.ValidatorAccumulatedCommission) ValidatorDistInfo {
	return ValidatorDistInfo{
		OperatorAddress:     operatorAddr,
		SelfBondRewards:     rewards,
		ValidatorCommission: commission,
	}
}

func validatorInfoHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		validatorAddr, ok := checkValidatorAddressVar(writer, req)
		if !ok {
			return
		}
		//commission
		res, err := common.QueryValidatorCommission(cliCtx, types.ModuleName, validatorAddr)
		if err != nil {
			rest.WriteErrorRes(writer, types.ErrInternalServer(types.DefaultCodespace))
		}
		var commission types.ValidatorAccumulatedCommission
		cliCtx.Cdc.MustUnmarshalJSON(res, &commission)

		//self bonded rewards
		delAddr := validatorAddr
		resp, Err := common.QueryDelegationRewards(cliCtx, types.ModuleName, validatorAddr, delAddr)
		if Err != nil {
			rest.WriteErrorRes(writer, types.ErrInternalServer(types.DefaultCodespace))
		}
		var rewards sdk.DecCoin
		cliCtx.Cdc.MustUnmarshalJSON(resp, &rewards)

		rest.PostProcessResponseBare(writer, cliCtx, NewValidatorDistInfo(validatorAddr, rewards, commission))
	}
}