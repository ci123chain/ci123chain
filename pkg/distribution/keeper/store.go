package keeper

import (
	"errors"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	types "github.com/tanhuiya/ci123chain/pkg/distribution/types"
	"github.com/tanhuiya/ci123chain/pkg/fc"
	"strconv"
)

// keeper of the staking store
type DistrKeeper struct {
	storeKey            sdk.StoreKey
	cdc                 *codec.Codec
	FeeCollectionKeeper fc.FcKeeper
	ak                  account.AccountKeeper
}

var (
	ValidatorCurrentRewardsPrefix = []byte("val")
	ValidatorsInfoPrefix = []byte("vals")
	DisrtKey = "distr"
)

// create a new keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, fck fc.FcKeeper, ak account.AccountKeeper) DistrKeeper {
	keeper := DistrKeeper{
		storeKey:            key,
		cdc:                 cdc,
		FeeCollectionKeeper: fck,
		ak:                  ak,
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
func (d *DistrKeeper) SetProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr, rewards sdk.Coin, height int64) {

	key := getKey(val, height)
	store := d.GetPreFixStore(ctx, ctx.ChainID())
	b := d.cdc.MustMarshalBinaryLengthPrefixed(rewards)
	store.Set(GetValidatorCurrentRewardsKey(key), b)
}

func (d *DistrKeeper) GetProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr, height int64) (rewards sdk.Coin) {

	key := getKey(val, height)
	store := d.GetPreFixStore(ctx, ctx.ChainID())
	b := store.Get(GetValidatorCurrentRewardsKey(key))
	if b == nil {
		return sdk.NewCoin(sdk.NewInt(0))
	}
	d.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

func (d *DistrKeeper) DeleteProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr) {
	store := d.GetPreFixStore(ctx, ctx.ChainID())
	store.Delete(GetValidatorCurrentRewardsKey(val))
}

//validator
func (d *DistrKeeper) SetValidatorCurrentRewards(ctx sdk.Context, val sdk.AccAddr, rewards sdk.Coin, height int64) {

	key := getKey(val, height)
	store := d.GetPreFixStore(ctx, ctx.ChainID())
	b := d.cdc.MustMarshalBinaryLengthPrefixed(rewards)
	store.Set(GetValidatorCurrentRewardsKey(key), b)
}

func (d *DistrKeeper) GetValidatorCurrentRewards(ctx sdk.Context, val sdk.AccAddr, height int64) (rewards sdk.Coin) {

	key := getKey(val, height)
	store := d.GetPreFixStore(ctx, ctx.ChainID())
	b := store.Get(GetValidatorCurrentRewardsKey(key))
	if b == nil {
		return sdk.NewCoin(sdk.NewInt(0))
	}
	d.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

func (d *DistrKeeper) DeleteValidatorOldRewardsRecord(ctx sdk.Context, val sdk.AccAddr) {

	store := d.GetPreFixStore(ctx, ctx.ChainID())
	b := store.Get(GetValidatorCurrentRewardsKey(val))
	if b == nil {
		return
	}
	store.Delete(GetValidatorCurrentRewardsKey(val))
}

//query
func (d *DistrKeeper) GetValCurrentRewards(ctx sdk.Context, val sdk.AccAddr) (rewards sdk.Coin, err error) {

	store := d.GetPreFixStore(ctx, ctx.ChainID())
	b := store.Get(GetValidatorCurrentRewardsKey(val))
	if b == nil {
		return sdk.NewCoin(sdk.NewInt(0)), errors.New("no such information")
	}
	d.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
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
func (d *DistrKeeper) GetPreviousProposer(ctx sdk.Context) (proposer sdk.AccAddr){

	store := d.GetPreFixStore(ctx, ctx.ChainID())
	b := store.Get(types.ProposerKey)
	if b == nil {
		//panic("Previous proposer not set")
		return sdk.AccAddr{}
	}
	d.cdc.MustUnmarshalBinaryLengthPrefixed(b, &proposer)
	return
}


func (d *DistrKeeper) SetPreviousProposer(ctx sdk.Context, proposer sdk.AccAddr) {
	store := d.GetPreFixStore(ctx, ctx.ChainID())
	b := d.cdc.MustMarshalBinaryLengthPrefixed(proposer)
	store.Set(types.ProposerKey, b)
}

//validatorsInfo
func (d *DistrKeeper) SetValidatorsInfo(ctx sdk.Context, bytes []byte, height int64) {
	key := []byte(strconv.FormatInt(height, 10))
	store := d.GetPreFixStore(ctx, ctx.ChainID())
	store.Set(GetValidatorsInfoKey(key), bytes)
}

func (d *DistrKeeper) GetValidatorsInfo(ctx sdk.Context, height int64) []byte{
	key := []byte(strconv.FormatInt(height, 10))
	store := d.GetPreFixStore(ctx, ctx.ChainID())
	bz := store.Get(GetValidatorsInfoKey(key))
	return bz
}

func (d *DistrKeeper) DeleteValidatorsInfo(ctx sdk.Context, height int64) {
	store := d.GetPreFixStore(ctx, ctx.ChainID())
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

func (d *DistrKeeper) GetPreFixStore(ctx sdk.Context, prefix string) sdk.KVStore{
	store := ctx.KVStore(d.storeKey).Prefix([]byte(prefix))
	return store
}