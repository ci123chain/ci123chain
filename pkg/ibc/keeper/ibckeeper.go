package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"github.com/tanhuiya/ci123chain/pkg/supply"
	"strconv"
	"time"
)

const Priv = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgO+x/1pjgqImlzWe+
fQj0E0ml/ajNet3lqenPtyvEwB+hRANCAASbLWrcFumBm7tzZKpCiPl/gzmVm1GI
2vwHa6qRkVdEjMpLIL7weErc1C+/ww81NBRgDGyNxiHq6ndBUNHxv9M3
-----END PRIVATE KEY-----`
// address: 0x505A74675dc9C71eF3CB5DF309256952917E801e


type IBCKeeper struct {
	AccountKeeper account.AccountKeeper
	SupplyKeeper  supply.Keeper
	StoreKey 	 sdk.StoreKey
}

func NewIBCKeeper(key sdk.StoreKey, AccountKeeper account.AccountKeeper, SupplyKeeper supply.Keeper) IBCKeeper {
	return IBCKeeper{
		StoreKey:	key,
		SupplyKeeper: SupplyKeeper,
		AccountKeeper:AccountKeeper,
	}
}

// 获取一个 ibcmsg
func (k IBCKeeper) GetFirstReadyIBCMsg(ctx sdk.Context) *types.IBCInfo {
	store := k.getStore(ctx)
	itr := sdk.KVStorePrefixIterator(store, []byte(types.StateKey + types.StateReady + types.TimestampKey))
	defer itr.Close()
	var ibc_msg *types.IBCInfo
	for {
		if !itr.Valid() {
			break
		}
		uniqueID := itr.Value()
		ibc_msg = k.GetIBCByUniqueID(ctx, uniqueID)
		break
	}

	// if ibcmsg is nil, search ibcmsg state = processing & timeout
	if ibc_msg == nil {
		preKey := types.StateKey + types.StateProcessing + types.TimestampKey
		itr = sdk.KVStorePrefixIterator(store, []byte(preKey))
		for {
			if !itr.Valid() {
				break
			}
			uniqueID := itr.Value()
			item := k.GetIBCByUniqueID(ctx, uniqueID)
			if time.Now().Unix() - item.ApplyTime.Unix() > types.TimeoutProcessing {
				// timeout
				ibc_msg = item
				break
			}
		}
	}
	return ibc_msg
}

// 申请处理某笔交易
func (k IBCKeeper) ApplyIBCMsg(ctx sdk.Context, uniqueID []byte, observerID []byte) (*types.ApplyReceipt, error) {
	ibcMsg := k.GetIBCByUniqueID(ctx, uniqueID)
	if ibcMsg == nil {
		return nil, errors.New(fmt.Sprintf("ibc tx not found with uniqueID = %s", string(uniqueID)))
	}
	if !ibcMsg.CanProcess() {
		return nil, errors.New(fmt.Sprintf("ibc tx not avaliable with uniqueID = %s, state = %s", string(uniqueID), ibcMsg.State))
	}
	// 修改处理人状态，以及时间
	ibcMsg.ApplyTime = time.Now()
	ibcMsg.ObserverID = observerID
	bankAddr, err := getBankAddress()
	if err != nil {
		return nil, err
	}
	ibcMsg.BankAddress = bankAddr

	// 以及抵押
	if ibcMsg.State == types.StateReady {
		// 抵押 给 module
		coinFrom := k.AccountKeeper.GetBalance(ctx, ibcMsg.FromAddress)
		_, valid := coinFrom.SafeSub(ibcMsg.Amount)
		if !valid {
			ibcMsg.State = types.StateCancel
			err = k.SetIBCMsg(ctx, *ibcMsg)
			return nil, errors.New("Infficient balance of account: " + ibcMsg.FromAddress.Hex())
		}
		if err1 := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, ibcMsg.FromAddress, types.ModuleName, ibcMsg.Amount); err1 != nil {
			ibcMsg.State = types.StateCancel
			err = k.SetIBCMsg(ctx, *ibcMsg)
			return nil, errors.New(err1.Error())
		}
	}

	// 修改状态
	ibcMsg.State = types.StateProcessing

	ibcMsgJson, err := json.Marshal(*ibcMsg)
	if err != nil {
		return nil, err
	}
	signedIbcMsg := types.ApplyReceipt{
		IBCMsgBytes:  ibcMsgJson,
	}
	priv, err := getPrivateKey()
	if err != nil {
		return nil, err
	}
	signedIbcMsg, err = signedIbcMsg.Sign(priv)
	if err != nil {
		return nil, err
	}

	// 保存状态
	err = k.SetIBCMsg(ctx, *ibcMsg)
	if err != nil {
		return nil, err
	}
	return &signedIbcMsg, nil
}

// 根据 uniqueID 获取 消息
func (k IBCKeeper) GetIBCByUniqueID(ctx sdk.Context,uniqueID []byte) *types.IBCInfo {
	store := k.getStore(ctx)
	bz := store.Get(uniqueID)
	if len(bz) < 1 {
		return nil
	}

	var ibcMsg types.IBCInfo
	err := types.IbcCdc.UnmarshalBinaryLengthPrefixed(bz, &ibcMsg)
	if err != nil {
		panic(err)
	}
	return &ibcMsg
}

// 保存 ibcmsg
func (k IBCKeeper) SetIBCMsg(ctx sdk.Context,ibcMsg types.IBCInfo) error {

	store := k.getStore(ctx)
	// 删除旧的索引
	oldIbcMsg := k.GetIBCByUniqueID(ctx, ibcMsg.UniqueID)
	if oldIbcMsg != nil {
		idxkey := types.StateKey + oldIbcMsg.State + types.TimestampKey + strconv.FormatInt(oldIbcMsg.ApplyTime.Unix(), 10) + types.UniqueKey + string(oldIbcMsg.UniqueID)
		store.Delete([]byte(idxkey))
	}

	bz, err := types.IbcCdc.MarshalBinaryLengthPrefixed(ibcMsg)
	if err != nil {
		return err
	}

	store.Set(ibcMsg.UniqueID, bz)

	// 保存索引结构
	//uniqueID := string(ibcMsg.UniqueID)

	idxkey := types.StateKey + ibcMsg.State + types.TimestampKey + strconv.FormatInt(ibcMsg.ApplyTime.Unix(), 10) + types.UniqueKey + string(ibcMsg.UniqueID)
	Logger(ctx).Debug(idxkey)
	store.Set([]byte(idxkey), ibcMsg.UniqueID)
	return nil
}

// 银行 转账到个人账户
func (k IBCKeeper) BankSend(ctx sdk.Context, ibcMsg types.IBCInfo) error {
	Logger(ctx).Debug(ibcMsg.BankAddress.Hex())
	err := k.AccountKeeper.Transfer(ctx, ibcMsg.BankAddress, ibcMsg.ToAddress, ibcMsg.Amount)
	if err != nil {
		return err
	}
	return nil
}

// 生成转账回执
func (k IBCKeeper) MakeBankReceipt(ctx sdk.Context, ibcMsg types.IBCInfo) (*types.BankReceipt, error) {
	bReceipt := types.NewBankReceipt(string(ibcMsg.UniqueID), string(ibcMsg.ObserverID))

	priv, err := getPrivateKey()
	if err != nil {
		return nil, err
	}
	bReceipt, err = bReceipt.Sign(priv)
	if err != nil {
		return nil, err
	}
	return bReceipt, nil
}

// 处理 回执
func (k IBCKeeper) ReceiveReceipt(ctx sdk.Context, receipt types.BankReceipt) (error) {
	uniqueID := receipt.UniqueID
	// todo uniqueID type
	ibcMsg := k.GetIBCByUniqueID(ctx, []byte(uniqueID))
	if ibcMsg == nil {
		return errors.New("IbcTx not found with uniqueID " + receipt.UniqueID)
	}
	if !bytes.Equal(ibcMsg.ObserverID, []byte(receipt.ObserverID)) {
		return errors.New("Got different observerID, expected same")
	}
	// 更新跨链交易状态
	ibcMsg.State = types.StateDone

	// 从 module 转给 目标
	macc := k.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	if macc == nil {
		return errors.New("Module account not set")
	}
	// 抵押 给 module
	coinFrom := k.AccountKeeper.GetBalance(ctx, macc.GetAddress())
	_, valid := coinFrom.SafeSub(ibcMsg.Amount)
	if !valid {
		return errors.New("Infficient balance of module account: " + macc.GetAddress().Hex())
	}
	if err1 := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ibcMsg.ToAddress, ibcMsg.Amount); err1 != nil {
		return errors.New(err1.Error())
	}

	err := k.SetIBCMsg(ctx, *ibcMsg)
	return err
}

func (k IBCKeeper) getStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.StoreKey)
}
