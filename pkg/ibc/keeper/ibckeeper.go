package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"github.com/tanhuiya/ci123chain/pkg/supply"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"time"
)

const Priv = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgp4qKKB0WCEfx7XiB
5Ul+GpjM1P5rqc6RhjD5OkTgl5OhRANCAATyFT0voXX7cA4PPtNstWleaTpwjvbS
J3+tMGTG67f+TdCfDxWYMpQYxLlE8VkbEzKWDwCYvDZRMKCQfv2ErNvb
-----END PRIVATE KEY-----`


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
	itr := sdk.KVStorePrefixIterator(store, []byte(types.StateKey + types.StateReady))
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

// 申请处理某笔交易
func (k IBCKeeper) ApplyIBCMsg(ctx sdk.Context, uniqueID []byte, observerID []byte) (*types.SignedIBCMsg, error) {
	ibcMsg := k.GetIBCByUniqueID(ctx, uniqueID)
	if ibcMsg == nil {
		return nil, errors.New(fmt.Sprintf("ibc tx not found with uniqueID = %s", hex.EncodeToString(uniqueID)))
	}
	if !ibcMsg.CanProcess() {
		return nil, errors.New(fmt.Sprintf("ibc tx not avaliable with uniqueID = %s, state = %s", hex.EncodeToString(uniqueID), ibcMsg.State))
	}
	// 修改处理人状态，以及时间
	ibcMsg.ApplyTime = time.Now()
	ibcMsg.ObserverID = observerID

	ibcMsgJson, err := types.IbcCdc.MarshalJSON(*ibcMsg)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(ibcMsgJson))

	signedIbcMsg := types.SignedIBCMsg{
		IBCMsgBytes:  ibcMsgJson,
	}
	priv, err := getPrivateKey()
	if err != nil {
		return nil, err
	}
	signedIbcMsg, err = signedIbcMsg.Sign(priv)
	return &signedIbcMsg, err
}

// 根据 uniqueID 获取 消息
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
	uniqueID := string(ibcMsg.UniqueID)
	idxkey := types.StateKey + types.StateReady + types.UniqueKey + uniqueID
	store.Set([]byte(idxkey), ibcMsg.UniqueID)
	return nil
}

func (k IBCKeeper) getStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.StoreKey)
}

func getPrivateKey() ([]byte, error) {
	priKey, err := cryptoutil.DecodePriv([]byte(Priv))
	if err != nil {
		return nil, err
	}
	priBz := cryptoutil.MarshalPrivateKey(priKey)
	return priBz, nil
}