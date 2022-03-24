package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)


// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {

		// Valsets
		case types.QueryCurrentValset:
			return queryCurrentValset(ctx, keeper)
		case types.QueryValsetRequest:
			return queryValsetRequestByNonce(ctx, path[1], keeper)
		//case types.QueryValsetConfirm:
		//	return queryValsetConfirm(ctx, path[1:], keeper)
		case types.QueryValsetConfirmsByNonce:
			return queryAllValsetConfirms(ctx, path[1], path[2], keeper)
		case types.QueryLastValsetRequests:
			return lastValsetRequests(ctx, keeper)
		case types.QueryLastPendingValsetRequestByAddr:
			return lastPendingValsetRequest(ctx, path[1], path[2], keeper)

		// Batches
		case types.QueryBatch:
			return queryBatch(ctx, path[1], path[2], path[3], keeper)
		case types.QueryBatchConfirms:
			return queryAllBatchConfirms(ctx, path[1], path[2], path[3], keeper)
		case types.QueryLastPendingBatchRequestByAddr:
			return lastPendingBatchRequest(ctx, path[1], path[2], keeper)
		case types.QueryLatestTxBatches:
			return lastBatchesRequest(ctx, path[1], keeper)
		//case types.QueryBatchFees:
		//	return queryBatchFees(ctx, keeper)


		// Token mappings
		case types.QueryDenomToERC20:
			return queryDenomToERC20(ctx, path[1], path[2], keeper)
		case types.QueryERC20ToDenom:
			return queryERC20ToDenom(ctx, path[1], path[2], keeper)
		case types.QueryDenomToERC721:
			return queryDenomToERC721(ctx, path[1], path[2], keeper)
		case types.QueryERC721ToDenom:
			return queryERC721ToDenom(ctx, path[1], path[2], keeper)

		// Pending transactions
		case types.QueryPendingSendToEths:
			return queryPendingSendToEth(ctx, path[1], path[2], path[3], keeper)

		// Event
		case types.QueryLastEventNonce:
			return queryLastEventNonce(ctx, path[1], path[2], keeper)

		case types.QueryLastValsetConfirmNonce:
			return queryLastValsetConfirmNonce(ctx, path[1], keeper)

		//case QueryTxId:
		//	return queryTxId(ctx, path[1], keeper)
		//case QueryEventNonce:
		//	return queryEventNonce(ctx, path[1], keeper)

		case types.QueryObservedEventNonce:
			return queryObservedEventNonce(ctx, path[1], keeper)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint", types.ModuleName)
		}
	}
}

func queryObservedEventNonce(ctx sdk.Context, gravityID string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID
	nonce := keeper.GetLastObservedEventNonceWithGid(ctx)
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, nonce)
	return buf, nil
}

func queryValsetRequestByNonce(ctx sdk.Context, nonceS string, keeper Keeper) ([]byte, error) {
	//keeper.currentGID = gravityID
	nonce, err := types.UInt64FromString(nonceS)
	if err != nil {
		return nil, err
	}

	valset := keeper.GetValset(ctx, nonce)
	if valset == nil {
		return nil, nil
	}
	// TODO: replace these with the GRPC response types
	// TODO: fix the use of module codec here
	res, err := codec.MarshalJSONIndent(types.GravityCodec, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// allValsetConfirmsByNonce returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func queryAllValsetConfirms(ctx sdk.Context, gravityID, nonceStr string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID
	nonce, err := types.UInt64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var confirms []*types.MsgValsetConfirm
	keeper.IterateValsetConfirmByNonceWithGID(ctx, nonce, func(_ []byte, c types.MsgValsetConfirm) bool {
		confirms = append(confirms, &c)
		return false
	})
	if len(confirms) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.GravityCodec, confirms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// allBatchConfirms returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func queryAllBatchConfirms(ctx sdk.Context, gravityID, nonceStr string, tokenContract string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID
	nonce, err := types.UInt64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var confirms []types.MsgConfirmBatch
	keeper.IterateBatchConfirmByNonceAndTokenContractWithGID(ctx, nonce, tokenContract, func(_ []byte, c types.MsgConfirmBatch) bool {
		confirms = append(confirms, c)
		return false
	})
	if len(confirms) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.GravityCodec, confirms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

const maxValsetRequestsReturned = 5

// lastValsetRequests returns up to maxValsetRequestsReturned valsets from the store
func lastValsetRequests(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var counter int
	var valReq []*types.Valset
	keeper.IterateValsets(ctx, func(_ []byte, val *types.Valset) bool {
		valReq = append(valReq, val)
		counter++
		return counter >= maxValsetRequestsReturned
	})
	if len(valReq) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.GravityCodec, valReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// lastPendingValsetRequest gets a list of validator sets that this validator has not signed
// limited by 100 sets per request.
func lastPendingValsetRequest(ctx sdk.Context, gravityID, operatorAddr string, keeper Keeper) ([]byte, error) {
	keeper.SetCurrentGid(gravityID)
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingValsetReq []*types.Valset
	keeper.IterateValsets(ctx, func(_ []byte, val *types.Valset) bool {
		// foundConfirm is true if the operatorAddr has signed the valset we are currently looking at
		foundConfirm := keeper.GetValsetConfirmByGID(ctx, val.Nonce, addr) != nil
		// if this valset has NOT been signed by operatorAddr, store it in pendingValsetReq
		// and exit the loop
		if !foundConfirm {
			pendingValsetReq = append(pendingValsetReq, val)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, TODO pagination
		if len(pendingValsetReq) > 100 {
			return true
		}
		// return false to continue the loop
		return false
	})
	if len(pendingValsetReq) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.GravityCodec, pendingValsetReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryCurrentValset(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	valset := keeper.GetCurrentValset(ctx)
	res, err := codec.MarshalJSONIndent(types.GravityCodec, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryValsetConfirm returns the confirm msg for single orchestrator address and nonce
// When nothing found a nil value is returned
//func queryValsetConfirm(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
//	nonce, err := types.UInt64FromString(path[0])
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
//	}
//
//	accAddress, err := sdk.AccAddressFromBech32(path[1])
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
//	}
//
//	valset := keeper.GetValsetConfirmByGID(ctx, nonce, accAddress)
//	if valset == nil {
//		return nil, nil
//	}
//	res, err := codec.MarshalJSONIndent(types.GravityCodec, *valset)
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
//	}
//
//	return res, nil
//}

type MultiSigUpdateResponse struct {
	Valset     types.Valset `json:"valset"`
	Signatures [][]byte     `json:"signatures,omitempty"`
}

// lastPendingBatchRequest gets the latest batch that has NOT been signed by operatorAddr
func lastPendingBatchRequest(ctx sdk.Context, gravityID, operatorAddr string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID

	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingBatchReq *types.OutgoingTxBatch
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		foundConfirm := keeper.GetBatchConfirmWithGID(ctx, batch.BatchNonce, batch.TokenContract, addr) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})
	if pendingBatchReq == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.GravityCodec, pendingBatchReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

const MaxResults = 100 // todo: impl pagination

// Gets MaxResults batches from store. Does not select by token type or anything
func lastBatchesRequest(ctx sdk.Context, gravityID string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID

	var batches []*types.OutgoingTxBatch
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return len(batches) == MaxResults
	})
	if len(batches) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.GravityCodec, batches)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryBatchFees(ctx sdk.Context, gravityID string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID

	val := types.QueryBatchFeeResponse{BatchFees: keeper.GetAllBatchFees(ctx)}
	res, err := codec.MarshalJSONIndent(types.GravityCodec, val)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


// queryBatch gets a batch by tokenContract and nonce
func queryBatch(ctx sdk.Context, gravityID string, nonce string, tokenContract string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID

	parsedNonce, err := types.UInt64FromString(nonce)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	if types.ValidateEthAddress(tokenContract) != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	foundBatch := keeper.GetOutgoingTXBatch(ctx, tokenContract, parsedNonce)
	if foundBatch == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not find tx batch")
	}
	res, err := codec.MarshalJSONIndent(types.GravityCodec, foundBatch)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	return res, nil
}


func queryDenomToERC20(ctx sdk.Context, gravityID string, denom string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID
	cosmos_originated, erc20, err := keeper.DenomToERC20Lookup(ctx, denom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	if !cosmos_originated {
		return nil, sdkerrors.Wrap(types.ErrQueryERC20, "erc20 not found")
	}
	var response types.QueryDenomToERC20Response
	response.CosmosOriginated = cosmos_originated
	response.Erc20 = erc20
	bytes, err := types.GravityCodec.MarshalJSON(response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return bytes, nil
	}
}

func queryERC20ToDenom(ctx sdk.Context, gravityID string, ERC20 string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID
	cosmos_originated, denom := keeper.ERC20ToDenomLookup(ctx, ERC20)
	if !cosmos_originated {
		return nil, sdkerrors.Wrap(types.ErrQueryDenom, "denom not found")
	}
	var response types.QueryERC20ToDenomResponse
	response.CosmosOriginated = cosmos_originated
	response.Denom = denom
	bytes, err := types.GravityCodec.MarshalJSON(response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return bytes, nil
	}
}

func queryDenomToERC721(ctx sdk.Context, gravityID string, denom string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID
	cosmos_originated, erc721, err := keeper.DenomToERC721Lookup(ctx, denom)
	if !cosmos_originated {
		return nil, sdkerrors.Wrap(types.ErrQueryDenom, "721 not found")
	}
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	var response types.QueryDenomToERC721Response
	response.CosmosOriginated = cosmos_originated
	response.Erc721 = erc721
	bytes, err := types.GravityCodec.MarshalJSON(response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return bytes, nil
	}
}

func queryERC721ToDenom(ctx sdk.Context, gravityID string, ERC20 string, keeper Keeper) ([]byte, error) {
	keeper.currentGID = gravityID
	cosmos_originated, denom := keeper.ERC721ToDenomLookup(ctx, ERC20)
	if !cosmos_originated {
		return nil, sdkerrors.Wrap(types.ErrQueryDenom, "denom721 not found")
	}
	var response types.QueryERC721ToDenomResponse
	response.CosmosOriginated = cosmos_originated
	response.Denom = denom
	bytes, err := types.GravityCodec.MarshalJSON(response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return bytes, nil
	}
}

func queryPendingSendToEth(ctx sdk.Context, gravityID string, senderAddr,  wlkContract string, k Keeper) ([]byte, error) {
	k.currentGID = gravityID
	//batches := k.GetOutgoingTxBatches(ctx)
	unbatched_tx := k.GetPoolTransactions(ctx)
	sender_address := senderAddr
	wlk_contract_address := wlkContract
	var eth_contract_address string
	if len(wlkContract) != 0 {
		contract_address, exists := k.GetMapedEthToken(ctx, wlk_contract_address)
		if !exists {
			return nil, fmt.Errorf("denom not a default coin: %s, and also not a ERC20 index", wlk_contract_address)
		}
		eth_contract_address = contract_address
	}

	res := types.QueryPendingSendToEthResponse{}
	//for _, batch := range batches {
	//	for _, tx := range batch.Transactions {
	//		if tx.Sender == sender_address {
	//			res.TransfersInBatches = append(res.TransfersInBatches, tx)
	//		}
	//	}
	//}
	for _, tx := range unbatched_tx {
		if len(sender_address) != 0 && len(eth_contract_address) != 0 {
			if (tx.Sender == sender_address) && (tx.Erc20Token.Contract == eth_contract_address) {
				res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
			}
		} else if len(sender_address) != 0 {
			if tx.Sender == sender_address {
				res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
			}
		} else if len(eth_contract_address) != 0 {
			if tx.Erc20Token.Contract == eth_contract_address {
				res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
			}
		} else {
			res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
		}
	}
	bytes, err := codec.MarshalJSONIndent(types.GravityCodec, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return bytes, nil
	}
}

func queryLastEventNonce(ctx sdk.Context, gravityID string, address string, k Keeper) ([]byte, error) {
	k.currentGID = gravityID
	addr := sdk.HexToAddress(address)
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, addr)
	x := types.UInt64Bytes(lastEventNonce)
	return x, nil
}

func queryLastValsetConfirmNonce(ctx sdk.Context, gravityID string, k Keeper) ([]byte, error) {
	k.currentGID = gravityID
	lastValsetConfirmNonce := k.GetLastValsetConfirmNonce(ctx)
	x := types.UInt64Bytes(lastValsetConfirmNonce)
	return x, nil
}


func queryTxId(ctx sdk.Context, txIdStr string, k Keeper) ([]byte, error) {
	txId, err := strconv.ParseUint(txIdStr, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error())
	}

	by, err := k.getTxIdState(ctx, txId)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "tx_id")
	}
	return by, nil
}

func queryEventNonce(ctx sdk.Context, eventNonceStr string, k Keeper) ([]byte, error) {
	laseNonce, err := strconv.ParseUint(eventNonceStr, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error())
	}

	by, err := k.getEventNonceState(ctx, laseNonce)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "event_nonce")
	}
	return by, nil
}