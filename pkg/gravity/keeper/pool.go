package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

// AddToOutgoingPool
// - checks a counterpart denominator exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, counterpartReceiver string, amount sdk.Coin, fee sdk.Coin, tokenType uint64) (uint64, error) {
	// If the coin is a gravity voucher, burn the coins. If not, check if there is a deployed ERC20 contract representing it.
	// If there is, lock the coins.
	var counterPartContract, wlkContract string
	if tokenType == ERC20 {
		_, counter, err := k.DenomToERC20Lookup(ctx, amount.Denom)
		if err != nil {
			return 0, err
		}
		contract, exist := k.GetMapedWlkToken(ctx, counter)
		if !exist {
			return 0, types.ErrMappedContractNotFound
		}
		wlkContract = contract
		counterPartContract = counter
		if k.IsWlkToken(wlkContract) {
			// lock coins in module
			if err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{amount}); err != nil {
				return 0, err
			}
		} else {
			// If it is an ethereum-originated asset we burn it
			// send coins to module in prep for burn
			if err := k.SupplyKeeper.SendCoinsFromEVMAccountToModule(ctx, sender, types.ModuleName, sdk.HexToAddress(wlkContract), amount.Amount.BigInt()); err != nil {
				return 0, err
			}
			// burn vouchers to send them back to ETH
			//if err := k.SupplyKeeper.BurnEVMCoin(ctx, types.ModuleName, sdk.HexToAddress(wlkContract), totalAmount.Amount.BigInt()); err != nil {
			//	panic(err)
			//}
		}
	} else if tokenType == ERC721 {
		_, counter, err := k.DenomToERC721Lookup(ctx, amount.Denom)
		if err != nil {
			return 0, err
		}
		contract, exist := k.GetMapedWRC721Token(ctx, counter)
		if !exist {
			return 0, types.ErrMappedContractNotFound
		}
		wlkContract = contract
		counterPartContract = counter
		if err := k.SupplyKeeper.Send721CoinsFromEVMAccountToModule(ctx, sender, types.ModuleName, sdk.HexToAddress(wlkContract), amount.Amount.BigInt()); err != nil {
			return 0, err
		}
	}

	//fee
	if err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{fee}); err != nil {
		return 0, err
	}
	// get next tx id from keeper
	nextID := k.autoIncrementID(ctx, types.KeyLastTXPoolID)
	// add relation to txid and tokentype
	k.SetTxidTokenType(ctx, nextID, tokenType)

	erc20Fee := types.NewSDKIntERC20Token(fee.Amount, counterPartContract)

	// construct outgoing tx, as part of this process we represent
	// the token as an ERC20 token since it is preparing to go to ETH
	// rather than the denom that is the input to this function.
	outgoing := &types.OutgoingTransferTx{
		Id:          nextID,
		Sender:      sender.String(),
		DestAddress: counterpartReceiver,
		Erc20Token:  types.NewSDKIntERC20Token(amount.Amount, counterPartContract),
		Erc20Fee:    erc20Fee,
	}

	// set the outgoing tx in the pool index
	if err := k.setPoolEntry(ctx, outgoing); err != nil {
		return 0, err
	}

	k.setTxIdState(ctx, nextID, txIdStatePending)

	// add a second index with the fee
	k.appendToUnbatchedTXIndex(ctx, counterPartContract, *erc20Fee, nextID)

	// todo: add second index for sender so that we can easily query: give pending Tx by sender
	// todo: what about a second index for receiver?

	poolEvent := sdk.NewEvent(
		types.EventTypeBridgeWithdrawalReceived,
		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
		sdk.NewAttribute([]byte(types.AttributeKeyBridgeChainID), []byte(k.currentGID)),
		sdk.NewAttribute([]byte(types.AttributeKeyOutgoingTXID), []byte(strconv.Itoa(int(nextID)))),
		sdk.NewAttribute([]byte(types.AttributeKeyNonce), []byte(fmt.Sprint(nextID))),
	)
	ctx.EventManager().EmitEvent(poolEvent)

	return nextID, nil
}

// RemoveFromOutgoingPoolAndRefund
// - checks that the provided tx actually exists
// - deletes the unbatched tx from the pool
// - issues the tokens back to the sender
func (k Keeper) RemoveFromOutgoingPoolAndRefund(ctx sdk.Context, txId uint64, sender sdk.AccAddress) error {
	// check that we actually have a tx with that id and what it's details are
	tx, err := k.getPoolEntry(ctx, txId)
	if err != nil {
		return err
	}

	found := false
	poolTx := k.GetPoolTransactions(ctx)
	for _, pTx := range poolTx {
		if pTx.Id == txId {
			found = true
		}
	}
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalid, "Id %d is in a batch", txId)
	}

	if strings.ToLower(tx.Sender) != strings.ToLower(sender.String()) {
		return sdkerrors.Wrapf(types.ErrInvalid, "Msg.Sender not Tx.Sender %d", txId)
	}

	// An inconsistent entry should never enter the store, but this is the ideal place to exploit
	// it such a bug if it did ever occur, so we should double check to be really sure
	if tx.Erc20Fee.Contract != tx.Erc20Token.Contract {
		return sdkerrors.Wrapf(types.ErrInvalid, "Inconsistent tokens to cancel!: %s %s", tx.Erc20Fee.Contract, tx.Erc20Token.Contract)
	}

	// delete this tx from both indexes
	k.removePoolEntry(ctx, txId)
	k.removeFromUnbatchedTXIndex(ctx, *tx.Erc20Fee, txId)

	k.setTxIdState(ctx, txId, txIDStateCancel)

	feeToRefundCoins := sdk.NewCoins(sdk.NewChainCoin(tx.Erc20Fee.Amount))
	if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, feeToRefundCoins); err != nil {
		return err
	}

	// reissue the amount and the fee
	tokenType := k.GetTxidTokenType(ctx, txId)
	if tokenType == ERC20 {
		wlkToken, _ := k.GetMapedWlkToken(ctx, tx.Erc20Token.Contract)
		if k.IsWlkToken(wlkToken) {
			if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(sdk.NewChainCoin(tx.Erc20Token.Amount))); err != nil {
				return err
			}
		} else {
			if err = k.SupplyKeeper.SendCoinsFromModuleToEVMAccount(ctx, sender, types.ModuleName, sdk.HexToAddress(wlkToken), tx.Erc20Token.Amount.BigInt()); err != nil {
				return sdkerrors.Wrap(err, "module transfer failed")
			}
		}
	} else {
		wlkToken, _ := k.GetMapedWRC721Token(ctx, tx.Erc20Token.Contract)
		if err = k.SupplyKeeper.Send721CoinsFromModuleToEVMAccount(ctx, sender, types.ModuleName, sdk.HexToAddress(wlkToken), tx.Erc20Token.Amount.BigInt()); err != nil {
			return sdkerrors.Wrap(err, "module transfer failed")
		}
	}

	poolEvent := sdk.NewEvent(
		types.EventTypeBridgeWithdrawCanceled,
		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
		sdk.NewAttribute([]byte(types.AttributeKeyBridgeChainID), []byte(k.currentGID)),
	)
	ctx.EventManager().EmitEvent(poolEvent)

	return nil
}

// appendToUnbatchedTXIndex add at the end when tx with same fee exists
func (k Keeper) appendToUnbatchedTXIndex(ctx sdk.Context, tokenContract string, fee types.ERC20Token, txID uint64) {
	store := k.getGidStore(ctx)
	idxKey := types.GetFeeSecondIndexKey(fee)
	var idSet types.IDSet
	if store.Has(idxKey) {
		bz := store.Get(idxKey)
		k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
	}
	idSet.Ids = append(idSet.Ids, txID)
	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
}

// appendToUnbatchedTXIndex add at the top when tx with same fee exists
func (k Keeper) prependToUnbatchedTXIndex(ctx sdk.Context, tokenContract string, fee types.ERC20Token, txID uint64) {
	store := k.getGidStore(ctx)
	idxKey := types.GetFeeSecondIndexKey(fee)
	var idSet types.IDSet
	if store.Has(idxKey) {
		bz := store.Get(idxKey)
		k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
	}
	idSet.Ids = append([]uint64{txID}, idSet.Ids...)
	store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
}

// removeFromUnbatchedTXIndex removes the tx from the index and makes it implicit no available anymore
func (k Keeper) removeFromUnbatchedTXIndex(ctx sdk.Context, fee types.ERC20Token, txID uint64) error {
	store := k.getGidStore(ctx)
	idxKey := types.GetFeeSecondIndexKey(fee)
	var idSet types.IDSet
	bz := store.Get(idxKey)
	if bz == nil {
		return sdkerrors.Wrap(types.ErrUnknown, "fee")
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &idSet)
	for i := range idSet.Ids {
		if idSet.Ids[i] == txID {
			idSet.Ids = append(idSet.Ids[0:i], idSet.Ids[i+1:]...)
			if len(idSet.Ids) != 0 {
				store.Set(idxKey, k.cdc.MustMarshalBinaryBare(&idSet))
			} else {
				store.Delete(idxKey)
			}
			return nil
		}
	}
	return sdkerrors.Wrap(types.ErrUnknown, "tx id")
}

func (k Keeper) setPoolEntry(ctx sdk.Context, val *types.OutgoingTransferTx) error {
	bz, err := k.cdc.MarshalBinaryBare(val)
	if err != nil {
		return err
	}
	store := k.getGidStore(ctx)
	store.Set(types.GetOutgoingTxPoolKey(val.Id), bz)
	return nil
}

func (k Keeper) getPoolEntry(ctx sdk.Context, id uint64) (*types.OutgoingTransferTx, error) {
	store := k.getGidStore(ctx)
	bz := store.Get(types.GetOutgoingTxPoolKey(id))
	if bz == nil {
		return nil, types.ErrUnknown.Wrap("Can not find this txid")
	}
	var r types.OutgoingTransferTx
	k.cdc.UnmarshalBinaryBare(bz, &r)
	return &r, nil
}

func (k Keeper) removePoolEntry(ctx sdk.Context, id uint64) {
	store := k.getGidStore(ctx)
	store.Delete(types.GetOutgoingTxPoolKey(id))
}

var (
	txIdStatePending = []byte("pending")
	txIdStateDone    = []byte("done")
	txIDStateCancel  = []byte("cancel")
)

func (k Keeper) setTxIdState(ctx sdk.Context, txId uint64, state []byte) {
	store := k.getGidStore(ctx)
	store.Set(types.GetTxIdKey(txId), state)
}

func (k Keeper) getTxIdState(ctx sdk.Context, txId uint64) ([]byte, error) {
	store := k.getGidStore(ctx)
	by := store.Get(types.GetTxIdKey(txId))
	if by == nil {
		return nil, types.ErrUnknown
	}
	return by, nil
}

var (
	eventNonceStateDone = []byte("done")
)

func (k Keeper) setEventNonceState(ctx sdk.Context, eventNonce uint64, state []byte) {
	store := k.getGidStore(ctx)
	store.Set(types.GetEventNonceKey(eventNonce), state)
}

func (k Keeper) getEventNonceState(ctx sdk.Context, eventNonce uint64) ([]byte, error) {
	store := k.getGidStore(ctx)
	by := store.Get(types.GetEventNonceKey(eventNonce))
	if by == nil {
		return nil, types.ErrUnknown
	}
	return by, nil
}

// GetPoolTransactions, grabs all transactions from the tx pool, useful for queries or genesis save/load
func (k Keeper) GetPoolTransactions(ctx sdk.Context) []*types.OutgoingTransferTx {
	store := k.getGidStore(ctx)
	// we must use the second index key here because transactions are left in the store, but removed
	// from the tx sorting key, while in batches
	iter := store.ReverseIterator(prefixRange(types.SecondIndexOutgoingTXFeeKey))
	var ret []*types.OutgoingTransferTx
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ids types.IDSet
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &ids)
		for _, id := range ids.Ids {
			tx, err := k.getPoolEntry(ctx, id)
			if err != nil {
				panic("Invalid id in tx index!")
			}
			ret = append(ret, tx)
		}
	}
	return ret
}

// IterateOutgoingPoolByFee iterates over the outgoing pool which is sorted by fee
func (k Keeper) IterateOutgoingPoolByFee(ctx sdk.Context, contract string, cb func(uint64, *types.OutgoingTransferTx) bool) {
	prefixStore := store.NewPrefixStore(k.getGidStore(ctx), types.SecondIndexOutgoingTXFeeKey)
	iter := prefixStore.ReverseIterator(prefixRange([]byte(contract)))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ids types.IDSet
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &ids)
		// cb returns true to stop early
		for _, id := range ids.Ids {
			tx, err := k.getPoolEntry(ctx, id)
			if err != nil {
				panic("Invalid id in tx index!")
			}
			if cb(id, tx) {
				return
			}
		}
	}
}

// GetBatchFeesByTokenType gets the fees the next batch of a given token type would
// have if created. This info is both presented to relayers for the purpose of determining
// when to request batches and also used by the batch creation process to decide not to create
// a new batch
func (k Keeper) GetBatchFeesByTokenType(ctx sdk.Context, tokenContractAddr string) *types.BatchFees {
	batchFeesMap := k.createBatchFees(ctx)
	return batchFeesMap[tokenContractAddr]
}

// GetAllBatchFees creates a fee entry for every batch type currently in the store
// this can be used by relayers to determine what batch types are desireable to request
func (k Keeper) GetAllBatchFees(ctx sdk.Context) (batchFees []*types.BatchFees) {
	batchFeesMap := k.createBatchFees(ctx)
	// create array of batchFees
	for _, batchFee := range batchFeesMap {
		// newBatchFee := types.BatchFees{
		// 	Token:         batchFee.Token,
		// 	TopOneHundred: batchFee.TopOneHundred,
		// }
		batchFees = append(batchFees, batchFee)
	}

	// quick sort by token to make this function safe for use
	// in consensus computations
	sort.Slice(batchFees, func(i, j int) bool {
		return batchFees[i].Token < batchFees[j].Token
	})

	return batchFees
}

// CreateBatchFees iterates over the outgoing pool and creates batch token fee map
func (k Keeper) createBatchFees(ctx sdk.Context) map[string]*types.BatchFees {
	prefixStore := store.NewPrefixStore(k.getGidStore(ctx), types.SecondIndexOutgoingTXFeeKey)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	batchFeesMap := make(map[string]*types.BatchFees)
	txCountMap := make(map[string]int)

	for ; iter.Valid(); iter.Next() {
		var ids types.IDSet
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &ids)

		// create a map to store the token contract address and its total fee
		// Parse the iterator key to get contract address & fee
		// If len(ids.Ids) > 1, multiply fee amount with len(ids.Ids) and add it to total fee amount

		key := iter.Key()
		tokenContractBytes := key[:types.ETHContractAddressLen]
		tokenContractAddr := string(tokenContractBytes)

		feeAmountBytes := key[len(tokenContractBytes):]
		feeAmount := big.NewInt(0).SetBytes(feeAmountBytes)

		for i := 0; i < len(ids.Ids); i++ {
			if txCountMap[tokenContractAddr] >= OutgoingTxBatchSize {
				break
			} else {
				// add fee amount
				if _, ok := batchFeesMap[tokenContractAddr]; ok {
					batchFeesMap[tokenContractAddr].TotalFees = batchFeesMap[tokenContractAddr].TotalFees.Add(sdk.NewIntFromBigInt(feeAmount))
				} else {
					batchFeesMap[tokenContractAddr] = &types.BatchFees{
						Token:     tokenContractAddr,
						TotalFees: sdk.NewIntFromBigInt(feeAmount)}
				}

				txCountMap[tokenContractAddr] = txCountMap[tokenContractAddr] + 1
			}
		}
	}

	return batchFeesMap
}

func (k Keeper) autoIncrementID(ctx sdk.Context, idKey []byte) uint64 {
	store := k.getGidStore(ctx)
	bz := store.Get(idKey)
	var id uint64 = 1
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(idKey, bz)
	return id
}
