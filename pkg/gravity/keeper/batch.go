package keeper

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"

	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

const OutgoingTxBatchSize = 100

// BuildOutgoingTXBatch starts the following process chain:
// - find bridged denominator for given voucher type
// - determine if a an unexecuted batch is already waiting for this token type, if so confirm the new batch would
//   have a higher total fees. If not exit withtout creating a batch
// - select available transactions from the outgoing transaction pool sorted by fee desc
// - persist an outgoing batch object with an incrementing ID = nonce
// - emit an event
func (k Keeper) BuildOutgoingTXBatch(ctx sdk.Context, contractAddress string, maxElements int, tokenType uint64, requestor common.Address) (*types.OutgoingTxBatch, error) {
	if maxElements == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "max elements value")
	}

	lastBatch := k.GetLastOutgoingBatchByTokenType(ctx, contractAddress)

	// lastBatch may be nil if there are no existing batches, we only need
	// to perform this check if a previous batch exists
	if lastBatch != nil {
		// this traverses the current tx pool for this token type and determines what
		// fees a hypothetical batch would have if created
		currentFees := k.GetBatchFeesByTokenType(ctx, contractAddress)
		if currentFees == nil {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "error getting fees from tx pool")
		}

		lastFees := lastBatch.GetFees()
		if lastFees.GT(currentFees.TotalFees) {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "new batch would not be more profitable")
		}
	}

	selectedTx, err := k.pickUnbatchedTX(ctx, contractAddress, maxElements)
	if len(selectedTx) == 0 || err != nil {
		return nil, err
	}
	nextID := k.autoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	batch := &types.OutgoingTxBatch{
		BatchNonce:    nextID,
		BatchTimeout:  k.getBatchTimeoutHeight(ctx),
		Transactions:  selectedTx,
		TokenContract: contractAddress,
		TokenType: tokenType,
	}
	k.StoreBatch(ctx, batch, requestor)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
		sdk.NewAttribute([]byte(types.AttributeKeyContract), []byte(k.GetBridgeContractAddress(ctx))),
		sdk.NewAttribute([]byte(types.AttributeKeyBridgeChainID), []byte(strconv.Itoa(int(k.GetBridgeChainID(ctx))))),
		sdk.NewAttribute([]byte(types.AttributeKeyOutgoingBatchID), []byte(fmt.Sprint(nextID))),
		sdk.NewAttribute([]byte(types.AttributeKeyNonce), []byte(fmt.Sprint(nextID))),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return batch, nil
}

/// This gets the batch timeout height in Ethereum blocks.
func (k Keeper) getBatchTimeoutHeight(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	currentCosmosHeight := ctx.BlockHeight()
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values
	// are zero because no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	heights := k.GetLastObservedEthereumBlockHeight(ctx)
	if heights.CosmosBlockHeight == 0 || heights.EthereumBlockHeight == 0 {
		return 0
	}
	// we project how long it has been in milliseconds since the last Ethereum block height was observed
	projected_millis := (uint64(currentCosmosHeight) - heights.CosmosBlockHeight) * params.AverageBlockTime
	// we convert that projection into the current Ethereum height using the average Ethereum block time in millis
	projected_current_ethereum_height := (projected_millis / params.AverageEthereumBlockTime) + heights.EthereumBlockHeight
	// we convert our target time for block timeouts (lets say 12 hours) into a number of blocks to
	// place on top of our projection of the current Ethereum block height.
	blocks_to_add := params.TargetBatchTimeout / params.AverageEthereumBlockTime
	return projected_current_ethereum_height + blocks_to_add
}

// OutgoingTxBatchExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It frees all the transactions in the batch, then cancels all earlier batches
func (k Keeper) OutgoingTxBatchExecuted(ctx sdk.Context, tokenContract string, nonce uint64) error {
	b := k.GetOutgoingTXBatch(ctx, tokenContract, nonce)
	if b == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}
	r := k.GetRequestBatch(ctx, tokenContract, nonce)
	if r == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}
	totalFee := new(big.Int)
	for _, tx := range b.Transactions {
		totalFee = totalFee.Add(totalFee, tx.Erc20Fee.Amount.BigInt())
	}
	if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress{r.Requestor}, sdk.NewCoins(sdk.NewChainCoin(sdk.NewIntFromBigInt(totalFee)))); err != nil {
		return err
	}

	// cleanup outgoing TX pool
	for _, tx := range b.Transactions {
		k.removePoolEntry(ctx, tx.Id)
		k.setTxIdState(ctx, tx.Id, txIdStateDone)
	}

	// Iterate through remaining batches
	k.IterateOutgoingTXBatches(ctx, func(key []byte, iter_batch *types.OutgoingTxBatch) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		if iter_batch.BatchNonce < b.BatchNonce {
			k.CancelOutgoingTXBatch(ctx, tokenContract, iter_batch.BatchNonce)
		}
		return false
	})

	// Delete batch since it is finished
	k.DeleteBatch(ctx, *b)
	return nil
}

// StoreBatch stores a transaction batch
func (k Keeper) StoreBatch(ctx sdk.Context, batch *types.OutgoingTxBatch, requestor common.Address) {
	store := ctx.KVStore(k.storeKey)
	// set the current block height when storing the batch
	batch.Block = uint64(ctx.BlockHeight())
	requestBatch := &types.RequestBatch{
		Requestor: requestor,
	}
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshalBinaryBare(batch))
	requestKey := types.GetOutgoingTxRequestBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(requestKey, k.cdc.MustMarshalBinaryBare(requestBatch))
	blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block)
	store.Set(blockKey, k.cdc.MustMarshalBinaryBare(batch))
}

// StoreBatchUnsafe stores a transaction batch w/o setting the height
func (k Keeper) StoreBatchUnsafe(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)
	store.Set(key, k.cdc.MustMarshalBinaryBare(batch))

	blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block)
	store.Set(blockKey, k.cdc.MustMarshalBinaryBare(batch))
}

// GetRequestBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetRequestBatch(ctx sdk.Context, tokenContract string, nonce uint64) *types.RequestBatch {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxRequestBatchKey(tokenContract, nonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	var b types.RequestBatch
	k.cdc.MustUnmarshalBinaryBare(bz, &b)
	return &b
}

// DeleteBatch deletes an outgoing transaction batch
func (k Keeper) DeleteBatch(ctx sdk.Context, batch types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce))
	store.Delete(types.GetOutgoingTxBatchBlockKey(batch.Block))
}

// pickUnbatchedTX find TX in pool and remove from "available" second index
func (k Keeper) pickUnbatchedTX(ctx sdk.Context, contractAddress string, maxElements int) ([]*types.OutgoingTransferTx, error) {
	var selectedTx []*types.OutgoingTransferTx
	var err error
	k.IterateOutgoingPoolByFee(ctx, contractAddress, func(txID uint64, tx *types.OutgoingTransferTx) bool {
		if tx != nil && tx.Erc20Fee != nil {
			selectedTx = append(selectedTx, tx)
			err = k.removeFromUnbatchedTXIndex(ctx, *tx.Erc20Fee, txID)
			return err != nil || len(selectedTx) == maxElements
		} else {
			// we found a nil, exit
			return true
		}
	})
	return selectedTx, err
}

// GetOutgoingTXBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTXBatch(ctx sdk.Context, tokenContract string, nonce uint64) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(tokenContract, nonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	var b types.OutgoingTxBatch
	k.cdc.MustUnmarshalBinaryBare(bz, &b)
	for _, tx := range b.Transactions {
		tx.Erc20Token.Contract = tokenContract
		tx.Erc20Fee.Contract = tokenContract
	}
	return &b
}

// CancelOutgoingTXBatch releases all TX in the batch and deletes the batch
func (k Keeper) CancelOutgoingTXBatch(ctx sdk.Context, tokenContract string, nonce uint64) error {
	batch := k.GetOutgoingTXBatch(ctx, tokenContract, nonce)
	if batch == nil {
		return types.ErrUnknown
	}
	for _, tx := range batch.Transactions {
		tx.Erc20Fee.Contract = tokenContract
		k.prependToUnbatchedTXIndex(ctx, tokenContract, *tx.Erc20Fee, tx.Id)
	}

	// Delete batch since it is finished
	k.DeleteBatch(ctx, *batch)

	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
		sdk.NewAttribute([]byte(types.AttributeKeyContract), []byte(k.GetBridgeContractAddress(ctx))),
		sdk.NewAttribute([]byte(types.AttributeKeyBridgeChainID), []byte(strconv.Itoa(int(k.GetBridgeChainID(ctx))))),
		sdk.NewAttribute([]byte(types.AttributeKeyOutgoingBatchID), []byte(fmt.Sprint(nonce))),
		sdk.NewAttribute([]byte(types.AttributeKeyNonce), []byte(fmt.Sprint(nonce))),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return nil
}

// IterateOutgoingTXBatches iterates through all outgoing batches in DESC order.
func (k Keeper) IterateOutgoingTXBatches(ctx sdk.Context, cb func(key []byte, batch *types.OutgoingTxBatch) bool) {
	prefixStore := store.NewPrefixStore(ctx.KVStore(k.storeKey), types.OutgoingTXBatchKey)
	iter := prefixStore.ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var batch types.OutgoingTxBatch
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &batch)
		// cb returns true to stop early
		if cb(iter.Key(), &batch) {
			break
		}
	}
}

// GetOutgoingTxBatches returns the outgoing tx batches
func (k Keeper) GetOutgoingTxBatches(ctx sdk.Context) (out []*types.OutgoingTxBatch) {
	k.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		out = append(out, batch)
		return false
	})
	return
}

// GetLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) GetLastOutgoingBatchByTokenType(ctx sdk.Context, token string) *types.OutgoingTxBatch {
	batches := k.GetOutgoingTxBatches(ctx)
	var lastBatch *types.OutgoingTxBatch = nil
	lastNonce := uint64(0)
	for _, batch := range batches {
		if batch.TokenContract == token && batch.BatchNonce > lastNonce {
			lastBatch = batch
			lastNonce = batch.BatchNonce
		}
	}
	return lastBatch
}

// SetLastSlashedBatchBlock sets the latest slashed Batch block height
func (k Keeper) SetLastSlashedBatchBlock(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedBatchBlock, types.UInt64Bytes(blockHeight))
}

// GetLastSlashedBatchBlock returns the latest slashed Batch block
func (k Keeper) GetLastSlashedBatchBlock(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastSlashedBatchBlock)
	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// GetUnSlashedBatches returns all the unslashed batches in state
func (k Keeper) GetUnSlashedBatches(ctx sdk.Context, maxHeight uint64) (out []*types.OutgoingTxBatch) {
	lastSlashedBatchBlock := k.GetLastSlashedBatchBlock(ctx)
	k.IterateBatchBySlashedBatchBlock(ctx, lastSlashedBatchBlock, maxHeight, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		if batch.Block > lastSlashedBatchBlock {
			out = append(out, batch)
		}
		return false
	})
	return
}

// IterateBatchBySlashedBatchBlock iterates through all Batch by last slashed Batch block in ASC order
func (k Keeper) IterateBatchBySlashedBatchBlock(ctx sdk.Context, lastSlashedBatchBlock uint64, maxHeight uint64, cb func([]byte, *types.OutgoingTxBatch) bool) {
	prefixStore := store.NewPrefixStore(ctx.KVStore(k.storeKey), types.OutgoingTXBatchBlockKey)
	iter := prefixStore.Iterator(types.UInt64Bytes(lastSlashedBatchBlock), types.UInt64Bytes(maxHeight))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var Batch types.OutgoingTxBatch
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &Batch)
		// cb returns true to stop early
		if cb(iter.Key(), &Batch) {
			break
		}
	}
}
