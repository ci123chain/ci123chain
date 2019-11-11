package keeper

import (
	"encoding/hex"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"github.com/tanhuiya/ci123chain/pkg/supply"
)

const (
	StateKey = ".state="
	UniqueKey = ".uniqueid="
	StateReady = "ready"
	StateProcessing = "processing"
	StateDone = "done"
)

type IBCKeeper struct {
	SupplyKeeper supply.Keeper
	StoreKey 	 sdk.StoreKey
}

func NewIBCKeeper(key sdk.StoreKey, supplyKeeper supply.Keeper) IBCKeeper {
	return IBCKeeper{
		StoreKey:	key,
		SupplyKeeper:supplyKeeper,
	}
}

// 获取一个 ibcmsg
func (k IBCKeeper) GetFirstReadyIBCMsg(ctx sdk.Context) *types.IBCMsg {
	store := k.getStore(ctx)
	itr := sdk.KVStorePrefixIterator(store, []byte(StateKey + StateReady))
	defer itr.Close()
	var ibc_msg *types.IBCMsg
	for {
		if !itr.Valid() {
			break
		}
		uniqueID := itr.Value()
		ibc_msg = k.GetIBCByUniqueID(ctx, uniqueID)
		break
	}
	return ibc_msg
}

func (k IBCKeeper) GetIBCByUniqueID(ctx sdk.Context,uniqueID []byte) *types.IBCMsg {
	store := k.getStore(ctx)
	bz := store.Get(uniqueID)
	if len(bz) < 1 {
		return nil
	}

	var ibcMsg types.IBCMsg
	err := types.IbcCdc.UnmarshalBinaryLengthPrefixed(bz, &ibcMsg)
	if err != nil {
		panic(err)
	}
	return &ibcMsg
}

// 保存 ibcmsg
func (k IBCKeeper) SetIBCMsg(ctx sdk.Context,ibcMsg types.IBCMsg) error {
	bz, err := types.IbcCdc.MarshalBinaryLengthPrefixed(ibcMsg)
	if err != nil {
		return err
	}
	store := k.getStore(ctx)
	store.Set(ibcMsg.UniqueID, bz)

	// 保存索引结构
	uniqueID := hex.EncodeToString(ibcMsg.UniqueID)
	idxkey := StateKey + StateReady + UniqueKey + uniqueID
	store.Set([]byte(idxkey), ibcMsg.UniqueID)
	return nil
}

func (k IBCKeeper) getStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.StoreKey)
}