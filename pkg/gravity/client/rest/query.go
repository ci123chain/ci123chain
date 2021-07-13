package rest

import (
	"fmt"
	"net/http"

	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/gorilla/mux"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

func getValsetRequestHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetRequest/%s", storeName, nonce), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorRes(w, "valset not found")
			return
		}

		var out types.Valset
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

// USED BY RUST
func batchByNonceHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]
		denom := vars[tokenAddress]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/batch/%s/%s", storeName, nonce, denom), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorRes(w, "valset not found")
			return
		}

		var out *types.OutgoingTxBatch
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

// USED BY RUST
func lastBatchesHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastBatches", storeName), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorRes(w, "valset not found")
			return
		}

		var out []*types.OutgoingTxBatch
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

// gets all the confirm messages for a given validator set nonce
func allValsetConfirmsHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetConfirms/%s", storeName, nonce), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorRes(w, "valset confirms not found")
			return
		}

		var out []*types.MsgValsetConfirm
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

// gets all the confirm messages for a given transaction batch
func allBatchConfirmsHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]
		denom := vars[tokenAddress]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/batchConfirms/%s/%s", storeName, nonce, denom), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorRes(w, "batch confirms not found")
			return
		}

		var out []types.MsgConfirmBatch
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

func lastValsetRequestsHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastValsetRequests", storeName), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorRes(w, "valset requests not found")
			return
		}

		var out []types.Valset
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

func lastValsetRequestsByAddressHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingValsetRequest/%s", storeName, operatorAddr), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorRes(w, "no pending valset requests found")
			return
		}

		var out []types.Valset
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

func lastBatchesByAddressHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingBatchRequest/%s", storeName, operatorAddr), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorRes(w, "no pending batch requests found")
			return
		}

		var out types.OutgoingTxBatch
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

func lastEventNonceByAddressHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastEventNonce/%s", storeName, operatorAddr), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}

		out := types.UInt64FromBytes(res)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

func lastLogicCallHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastLogicCalls", storeName), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), res)
	}
}

func currentValsetHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/currentValset", storeName), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		var out types.Valset
		cliCtx.Cdc.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

func denomToERC20Handler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[denom]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/DenomToERC20/%s", storeName, denom), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), res)
	}
}

func ERC20ToDenomHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ERC20 := vars[tokenAddress]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/ERC20ToDenom/%s", storeName, ERC20), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), res)
	}
}
