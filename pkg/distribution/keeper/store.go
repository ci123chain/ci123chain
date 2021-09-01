package keeper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	//"github.com/ci123chain/ci123chain/pkg/couchdb"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/params"
	staking "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/supply"
	dbm "github.com/tendermint/tm-db"
	"strconv"
)

// keeper of the staking store
type DistrKeeper struct {
	storeKey            sdk.StoreKey
	cdc                 *codec.Codec
	SupplyKeeper        supply.Keeper
	FeeCollectorName    string
	AccountKeeper       account.AccountKeeper
	ParamSpace          params.Subspace
	StakingKeeper       staking.StakingKeeper
	cdb					dbm.DB
}

var (
	ValidatorCurrentRewardsPrefix = []byte("val")
	ValidatorsInfoPrefix = []byte("vals")
	DisrtKey = "distr"
)

// create a new keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk supply.Keeper, ak account.AccountKeeper, feeCollector string, paramSpace params.Subspace,
	stakingKeeper staking.StakingKeeper, cdb dbm.DB) DistrKeeper {

	// ensure distribution module account is set
	if addr := sk.GetModuleAddress(types.ModuleName); addr.Bytes() == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	keeper := DistrKeeper{
		storeKey:            key,
		cdc:                 cdc,
		SupplyKeeper:        sk,
		AccountKeeper:       ak,
		FeeCollectorName:    feeCollector,
		ParamSpace:          paramSpace,
		StakingKeeper:       stakingKeeper,
		cdb:				 cdb,
	}
	return keeper
}

func GetValidatorCurrentRewardsKey(v sdk.AccAddr) []byte {
	return append(ValidatorCurrentRewardsPrefix, v...)
}

func GetValidatorsInfoKey(v []byte) []byte {
	return append(ValidatorsInfoPrefix, v...)
}

//proposer
func (k *DistrKeeper) SetProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr, rewards sdk.Coin, height int64) {

	key := getKey(val, height)
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(rewards)
	store.Set(GetValidatorCurrentRewardsKey(key), b)
}

func (k *DistrKeeper) GetProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr, height int64) (rewards sdk.Coin) {

	key := getKey(val, height)
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetValidatorCurrentRewardsKey(key))
	if b == nil {
		return sdk.NewChainCoin(sdk.NewInt(0))
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

func (k *DistrKeeper) DeleteProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetValidatorCurrentRewardsKey(val))
}


// delete current rewards for a validator
func (k *DistrKeeper) DeleteValidatorCurrentRewards(ctx sdk.Context, val sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorCurrentRewardsKey(val))
}

//validator
func (k *DistrKeeper) SetValidatorCurrentRewards(ctx sdk.Context, val sdk.AccAddress, rewards types.ValidatorCurrentRewards) {

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&rewards)
	store.Set(types.GetValidatorCurrentRewardsKey(val), b)
}

func (k *DistrKeeper) GetValidatorCurrentRewards(ctx sdk.Context, val sdk.AccAddress) (rewards types.ValidatorCurrentRewards) {

	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorCurrentRewardsKey(val))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

func (k *DistrKeeper) DeleteValidatorOldRewardsRecord(ctx sdk.Context, val sdk.AccAddr) {

	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetValidatorCurrentRewardsKey(val))
	if b == nil {
		return
	}
	store.Delete(GetValidatorCurrentRewardsKey(val))
}

//query
func (k *DistrKeeper) GetValCurrentRewards(ctx sdk.Context, val sdk.AccAddr) (rewards sdk.Coin, err error) {

	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetValidatorCurrentRewardsKey(val))
	if b == nil {
		return sdk.NewChainCoin(sdk.NewInt(0)), errors.New("no such information")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}
/*
func (d *DistrKeeper) DistributeRewardsToValidators(ctx sdk.Context, proposer sdk.AccAddress, fee sdk.Coin) {

	account := d.ak.GetAccount(ctx, proposer)
	accCoin := account.GetCoin()
	accCoin.SafeAdd(fee)
	d.ak.SetAccount(ctx, account)
	//n := float32(0.05)
	//mulNum := length + n
	//var v float32
	//var mulFee = float32(fee)
	//v = mulFee/mulNum
	//fmt.Print(v)
	//val := v * 0.05
	//value := types.Coin(uint64(val))
	//proposerAcc := d.ak.GetAccount(ctx, proposer)
	//accCoin := proposerAcc.GetCoin()
	//accCoin.SafeAdd(value)
	//d.ak.SetAccount(ctx, proposerAcc)
	//
	//validatorVal := types.Coin(uint64(v))
	//for i, _ := range validators {
	//	validatorAcc := d.ak.GetAccount(ctx, validators[i])
	//	accCoin := validatorAcc.GetCoin()
	//	accCoin.SafeAdd(validatorVal)
	//	d.ak.SetAccount(ctx, validatorAcc)
	//}
}
*/

//lastProposer
func (k *DistrKeeper) GetPreviousProposerAddr(ctx sdk.Context) (proposer sdk.AccAddr){

	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ProposerKey)
	if b == nil {
		//panic("Previous proposer not set")
		return sdk.AccAddr{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &proposer)
	return
}


func (k *DistrKeeper) SetPreviousProposerAddr(ctx sdk.Context, proposer sdk.AccAddr) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(proposer)
	store.Set(types.ProposerKey, b)
}

//validatorsInfo
func (k *DistrKeeper) SetValidatorsInfo(ctx sdk.Context, bytes []byte, height int64) {
	key := []byte(strconv.FormatInt(height, 10))
	store := ctx.KVStore(k.storeKey)
	store.Set(GetValidatorsInfoKey(key), bytes)
}

func (k *DistrKeeper) GetValidatorsInfo(ctx sdk.Context, height int64) []byte{
	key := []byte(strconv.FormatInt(height, 10))
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetValidatorsInfoKey(key))
	return bz
}

func (k *DistrKeeper) DeleteValidatorsInfo(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(strconv.FormatInt(height, 10))
	bz := store.Get(GetValidatorsInfoKey(key))
	if bz == nil {
		return
	}
	store.Delete(GetValidatorsInfoKey(key))
}

func getKey(val sdk.AccAddr, height int64) sdk.AccAddr {
	add := fmt.Sprintf("%X", val)
	h := strconv.FormatInt(height, 10)
	tKey := add + h
	key := sdk.AccAddr([]byte(tKey))
	return key
}

func (k *DistrKeeper) GetPreFixStore(ctx sdk.Context, prefix string) sdk.KVStore{
	store := ctx.KVStore(k.storeKey).Prefix([]byte(prefix))
	return store
}

func (k *DistrKeeper) GetFeePool(ctx sdk.Context) (feePool types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.FeePoolKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &feePool)
	return
}

func (k *DistrKeeper) SetFeePool(ctx sdk.Context, feePool types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&feePool)
	store.Set(types.FeePoolKey, b)
}

func (k *DistrKeeper) DeleteValidatorAccumulatedCommission(ctx sdk.Context, val sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorAccumulatedCommissionKey(val))
}

func (k *DistrKeeper) GetValidatorAccumulatedCommission(ctx sdk.Context, val sdk.AccAddress) (commission types.ValidatorAccumulatedCommission) {

	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorAccumulatedCommissionKey(val))
	if b == nil {
		return types.ValidatorAccumulatedCommission{
			Commission: sdk.DecCoin{
				Denom:  sdk.ChainCoinDenom,
				Amount: sdk.NewDec(0),
			},
		}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &commission)
	return
}

func (k *DistrKeeper) SetValidatorAccumulatedCommission(ctx sdk.Context, val sdk.AccAddress, commission types.ValidatorAccumulatedCommission) {
	var bz []byte

	store := ctx.KVStore(k.storeKey)
	if commission.Commission.IsZero() {
		bz = k.cdc.MustMarshalBinaryLengthPrefixed(&types.ValidatorAccumulatedCommission{})
	}else {
		bz = k.cdc.MustMarshalBinaryLengthPrefixed(&commission)
	}
	store.Set(types.GetValidatorAccumulatedCommissionKey(val), bz)
}

func (k *DistrKeeper) DeleteValidatorOutstandingRewards(ctx sdk.Context, val sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorOutstandingRewardsKey(val))
}

func (k *DistrKeeper) GetValidatorOutstandingRewards(ctx sdk.Context, val sdk.AccAddress) (rewards types.ValidatorOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorOutstandingRewardsKey(val))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rewards)
	return
}

func (k *DistrKeeper) SetValidatorOutstandingRewards(ctx sdk.Context, val sdk.AccAddress, rewards types.ValidatorOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&rewards)
	store.Set(types.GetValidatorOutstandingRewardsKey(val), b)
}

func (k DistrKeeper) GetDelegatorWithdrawAddr(ctx sdk.Context, delAddr sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDelegatorWithdrawAddrKey(delAddr))
	if b == nil {
		return delAddr
	}
	return sdk.ToAccAddress(b)
}

func (k DistrKeeper) SetDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDelegatorWithdrawAddrKey(delAddr), withdrawAddr.Bytes())
}

// delete a delegator withdraw addr
func (k DistrKeeper) DeleteDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegatorWithdrawAddrKey(delAddr))
}

// get historical rewards for a particular period
func (k DistrKeeper) GetValidatorHistoricalRewards(ctx sdk.Context, val sdk.AccAddress, period uint64) (rewards types.ValidatorHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorHistoricalRewardsKey(val, period))
	if b == nil {
		rewards = types.NewValidatorHistoricalRewards(sdk.NewEmptyDecCoin(), 1)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

// set historical rewards for a particular period
func (k DistrKeeper) SetValidatorHistoricalRewards(ctx sdk.Context, val sdk.AccAddress, period uint64, rewards types.ValidatorHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&rewards)
	store.Set(types.GetValidatorHistoricalRewardsKey(val, period), b)
}

// delete historical rewards for a validator
func (k DistrKeeper) DeleteValidatorHistoricalRewards(ctx sdk.Context, val sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetValidatorHistoricalRewardsPrefix(val)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStorePrefixIterator(store, prefix)
	}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		realKey := iterator.Key()
		store.Delete(realKey)
	}
}

// delete a historical reward
func (k DistrKeeper) DeleteValidatorHistoricalReward(ctx sdk.Context, val sdk.AccAddress, period uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorHistoricalRewardsKey(val, period))
}

// get the starting info associated with a delegator
func (k DistrKeeper) GetDelegatorStartingInfo(ctx sdk.Context, val sdk.AccAddress, del sdk.AccAddress) (period types.DelegatorStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDelegatorStartingInfoKey(val, del))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &period)
	return
}

func (k DistrKeeper) SetDelegatorStartingInfo(ctx sdk.Context, val, del sdk.AccAddress, period types.DelegatorStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&period)
	store.Set(types.GetDelegatorStartingInfoKey(val, del), b)
}

// check existence of the starting info associated with a delegator
func (k DistrKeeper) HasDelegatorStartingInfo(ctx sdk.Context, val sdk.AccAddress, del sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDelegatorStartingInfoKey(val, del))
}

// delete the starting info associated with a delegator
func (k DistrKeeper) DeleteDelegatorStartingInfo(ctx sdk.Context, val sdk.AccAddress, del sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegatorStartingInfoKey(val, del))
}

// get slash event for height
func (k DistrKeeper) GetValidatorSlashEvent(ctx sdk.Context, val sdk.AccAddress, height, period uint64) (event types.ValidatorSlashEvent, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorSlashEventKey(val, height, period))
	if b == nil {
		return types.ValidatorSlashEvent{}, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &event)
	return event, true
}

func (k DistrKeeper) SetValidatorSlashEvent(ctx sdk.Context, val sdk.AccAddress, height, period uint64, event types.ValidatorSlashEvent) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(&event)
	store.Set(types.GetValidatorSlashEventKey(val, height, period), b)
}


// delete slash events for a particular validator
func (k DistrKeeper) DeleteValidatorSlashEvents(ctx sdk.Context, val sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetValidatorSlashEventPrefix(val)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStorePrefixIterator(store, prefix)
	}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		realKey := iterator.Key()
		//_, ok := iterator.(*couchdb.CouchIterator)
		//if ok {
		//	realKey = sdk.GetRealKey(iterator.Key())
		//}
		store.Delete(realKey)
	}
}

// FundCommunityPool allows an account to directly fund the community fund pool.
// The amount is first added to the distribution module account and then directly
// added to the pool. An error is returned if the amount cannot be sent to the
// module account.
func (k DistrKeeper) FundCommunityPool(ctx sdk.Context, amount sdk.Coin, sender sdk.AccAddress) error {

	if err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return err
	}
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinFromCoin(amount))
	k.SetFeePool(ctx, feePool)
	return nil
}

// withdraw validator commission
func (k DistrKeeper) WithdrawValidatorCommission(ctx sdk.Context, valAddr sdk.AccAddress) (sdk.Coin, error) {
	// fetch validator accumulated commission
	accumCommission := k.GetValidatorAccumulatedCommission(ctx, valAddr)
	if accumCommission.Commission.IsZero() {
		return sdk.NewEmptyCoin(), types.ErrNoValidatorCommission
	}

	commission, remainder := accumCommission.Commission.TruncateDecimal()
	k.SetValidatorAccumulatedCommission(ctx, valAddr, types.ValidatorAccumulatedCommission{Commission: remainder}) // leave remainder to withdraw later

	// update outstanding
	outstanding := k.GetValidatorOutstandingRewards(ctx, valAddr).Rewards
	k.SetValidatorOutstandingRewards(ctx, valAddr, types.ValidatorOutstandingRewards{Rewards: outstanding.Sub(sdk.NewDecCoinFromCoin(commission))})

	if !commission.IsZero() {
		accAddr := valAddr
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, accAddr)
		err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, sdk.NewCoins(commission))
		if err != nil {
			return sdk.NewEmptyCoin(), err
		}
	}

	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeWithdrawCommission,
	//		sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
	//	),
	//)

	return commission, nil
}


// withdraw rewards from a delegation
func (k DistrKeeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) (sdk.Coin, error) {
	val := k.StakingKeeper.Validator(ctx, valAddr)
	if val == nil {
		return sdk.NewEmptyCoin(), types.ErrNoValidatorExist
	}

	del := k.StakingKeeper.Delegation(ctx, delAddr, valAddr)
	if del == nil {
		return sdk.NewEmptyCoin(), types.ErrNoDelegationExist
	}

	// withdraw rewards
	rewards, err := k.withdrawDelegationRewards(ctx, val, del)
	if err != nil {
		return sdk.NewEmptyCoin(), err
	}

	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeWithdrawRewards,
	//		sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
	//		sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
	//	),
	//)

	// reinitialize the delegation
	k.initializeDelegation(ctx, valAddr, delAddr)
	return rewards, nil
}

// SetWithdrawAddr sets a new address that will receive the rewards upon withdrawal
func (k DistrKeeper) SetWithdrawAddr(ctx sdk.Context, delegatorAddr sdk.AccAddress, withdrawAddr sdk.AccAddress) error {
	/*if k.blacklistedAddrs[withdrawAddr.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is blacklisted from receiving external funds", withdrawAddr)
	}*/

	if !k.GetWithdrawAddrEnabled(ctx) {
		return types.ErrSetWithdrawAddrDisabled
	}

	k.SetDelegatorWithdrawAddr(ctx, delegatorAddr, withdrawAddr)
	return nil
}

func (k DistrKeeper) IterateDelegatorWithdrawAddrs(ctx sdk.Context, handler func(del sdk.AccAddress, addr sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	//prefix := types.DelegatorWithdrawAddrPrefix
	//iter := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	iter := sdk.KVStorePrefixIterator(store, types.DelegatorWithdrawAddrPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.ToAccAddress(iter.Value())
		del := addr
		if handler(del, addr) {
			break
		}
	}
}

// GetValidatorSlashEventAddressHeight creates the height from a validator's slash event key.
func GetValidatorSlashEventAddressHeight(key []byte) (valAddr sdk.AccAddress, height uint64) {
	// key is in the format:
	// 0x08<valAddrLen (1 Byte)><valAddr_Bytes><height>: ValidatorSlashEvent
	valAddrLen := int(key[1])
	valAddr = sdk.ToAccAddress(key[2 : 2+valAddrLen])
	startB := 2 + valAddrLen
	b := key[startB : startB+8] // the next 8 bytes represent the height
	height = binary.BigEndian.Uint64(b)
	return
}

func (k DistrKeeper) IterateValidatorSlashEvents(ctx sdk.Context, handler func(val sdk.AccAddress, height uint64, event types.ValidatorSlashEvent) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSlashEventPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var event types.ValidatorSlashEvent
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &event)
		val, height := GetValidatorSlashEventAddressHeight(iter.Key())
		if handler(val, height, event) {
			break
		}
	}
}


func (k DistrKeeper) GetPreviousProposerConsAddr(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposerKey)
	if bz == nil {
		panic("previous proposer not set")
	}

	var addrValue sdk.AccAddr
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz[:], &addrValue)
	return sdk.ToAccAddress(addrValue)
}


// iterate validator outstanding rewards
func (k DistrKeeper) IterateValidatorOutstandingRewards(ctx sdk.Context, handler func(val sdk.AccAddress, rewards types.ValidatorOutstandingRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorOutstandingRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		rewards := types.ValidatorOutstandingRewards{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rewards)
		addr := sdk.ToAccAddress(iter.Key()[2:])
		if handler(addr, rewards) {
			break
		}
	}
}

// iterate over accumulated commissions
func (k DistrKeeper) IterateValidatorAccumulatedCommissions(ctx sdk.Context, handler func(val sdk.AccAddress, commission types.ValidatorAccumulatedCommission) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorAccumulatedCommissionPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var commission types.ValidatorAccumulatedCommission
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &commission)
		addr := sdk.ToAccAddress(iter.Key()[2:])
		if handler(addr, commission) {
			break
		}
	}
}

func GetValidatorHistoricalRewardsAddressPeriod(key []byte) (valAddr sdk.AccAddress, period uint64) {
	// key is in the format:
	// 0x05<valAddrLen (1 Byte)><valAddr_Bytes><period_Bytes>
	valAddrLen := int(key[1])
	valAddr = sdk.ToAccAddress(key[2 : 2+valAddrLen])
	b := key[2+valAddrLen:]
	if len(b) != 8 {
		panic("unexpected key length")
	}
	period = binary.LittleEndian.Uint64(b)
	return
}


// iterate over historical rewards
func (k DistrKeeper) IterateValidatorHistoricalRewards(ctx sdk.Context, handler func(val sdk.AccAddress, period uint64, rewards types.ValidatorHistoricalRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.ValidatorHistoricalRewards
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rewards)
		addr, period := GetValidatorHistoricalRewardsAddressPeriod(iter.Key())
		if handler(addr, period, rewards) {
			break
		}
	}
}


func (k DistrKeeper) IterateValidatorCurrentRewards(ctx sdk.Context, handler func(val sdk.AccAddress, rewards types.ValidatorCurrentRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorCurrentRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.ValidatorCurrentRewards
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rewards)
		addr := sdk.ToAccAddress(iter.Key()[2:])
		if handler(addr, rewards) {
			break
		}
	}
}

// GetDelegatorStartingInfoAddresses creates the addresses from a delegator starting info key.
func GetDelegatorStartingInfoAddresses(key []byte) (valAddr sdk.AccAddress, delAddr sdk.AccAddress) {
	// key is in the format:
	// 0x04<valAddrLen (1 Byte)><valAddr_Bytes><accAddrLen (1 Byte)><accAddr_Bytes>
	valAddrLen := int(key[1])
	valAddr = sdk.ToAccAddress(key[2 : 2+valAddrLen])
	delAddrLen := int(key[2+valAddrLen])
	delAddr = sdk.ToAccAddress(key[3+valAddrLen:])
	if len(delAddr.Bytes()) != delAddrLen {
		panic("unexpected key length")
	}

	return
}

// iterate over delegator starting infos
func (k DistrKeeper) IterateDelegatorStartingInfos(ctx sdk.Context, handler func(val sdk.AccAddress, del sdk.AccAddress, info types.DelegatorStartingInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DelegatorStartingInfoPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var info types.DelegatorStartingInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &info)
		val, del := GetDelegatorStartingInfoAddresses(iter.Key())
		if handler(val, del, info) {
			break
		}
	}
}
