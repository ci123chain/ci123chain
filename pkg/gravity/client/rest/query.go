package rest

import (
	"encoding/binary"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

func getValsetRequestByNonceHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
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
		gravityID := vars[gravity_id]
		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/batch/%s/%s/%s", storeName, gravityID, nonce, denom), nil, false)
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
		vars := mux.Vars(r)
		gravityID := vars[gravity_id]
		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastBatches/%s", storeName, gravityID), nil, false)
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
		gravityID := vars[gravity_id]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetConfirms/%s/%s", storeName, gravityID, nonce), nil, false)
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
		gravityID := vars[gravity_id]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/batchConfirms/%s/%s/%s", storeName, gravityID, nonce, denom), nil, false)
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
		gravityID := vars[gravity_id]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingValsetRequest/%s/%s", storeName, gravityID, operatorAddr), nil, false)
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
		gravityID := vars[gravity_id]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingBatchRequest/%s/%s", storeName, gravityID, operatorAddr), nil, false)
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
		gravityID := vars[gravity_id]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastEventNonce/%s/%s", storeName, gravityID, operatorAddr), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}

		out := types.UInt64FromBytes(res)
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
	}
}

//func lastValsetConfirmNonceHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastValsetConfirmNonce/%s", storeName), nil, false)
//		if err != nil {
//			rest.WriteErrorRes(w, err.Error())
//			return
//		}
//
//		out := types.UInt64FromBytes(res)
//		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), out)
//	}
//}


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
		gravityID := vars[gravity_id]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/DenomToERC20/%s/%s", storeName, gravityID, denom), nil, false)
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
		gravityID := vars[gravity_id]
		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/ERC20ToDenom/%s/%s", storeName, gravityID, ERC20), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), res)
	}
}

func denomToERC721Handler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[denom]
		gravityID := vars[gravity_id]
		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/DenomToERC721/%s/%s", storeName, gravityID, denom), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), res)
	}
}

func ERC721ToDenomHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ERC721 := vars[tokenAddress]
		gravityID := vars[gravity_id]
		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/ERC721ToDenom/%s/%s", storeName, gravityID, ERC721), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), res)
	}
}

func queryTxIdHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		txId := vars[txId]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/txId/%s", storeName, txId), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), string(res))
	}
}

func queryEventNonceHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		eventNonce := vars[eventNonce]

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/eventNonce/%s", storeName, eventNonce), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), string(res))
	}
}

func queryObservedEventNonceHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/observedEventNonce", storeName), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), int64(binary.BigEndian.Uint64(res)))
	}
}

func queryPendingSendToEthHandler(cliCtx context.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		sender := params.Get("sender")
		wlk_contract := params.Get("wlkContract")
		gravityID := params.Get(gravity_id)

		res, height, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/PendingSendToEth/%s/%s/%s", storeName, gravityID, sender, wlk_contract), nil, false)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx.WithHeight(height), string(res))
	}
}