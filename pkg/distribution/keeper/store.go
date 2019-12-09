package keeper


import (
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
	ValidatorCurrentRewardsPrefix        = []byte("val")
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

func (d *DistrKeeper) SetProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr, rewards sdk.Coin, height int64) {

	key := getKey(val, height)
	store := ctx.KVStore(d.storeKey)
	b := d.cdc.MustMarshalBinaryLengthPrefixed(rewards)
	store.Set(GetValidatorCurrentRewardsKey(key), b)
}

func (d *DistrKeeper) GetProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr, height int64) (rewards sdk.Coin) {

	key := getKey(val, height)
	store := ctx.KVStore(d.storeKey)
	b := store.Get(GetValidatorCurrentRewardsKey(key))
	if b == nil {
		return sdk.NewCoin()
	}
	d.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

func (d *DistrKeeper) GetProCurrentRewards(ctx sdk.Context, val sdk.AccAddr) (rewards sdk.Coin) {

	store := ctx.KVStore(d.storeKey)
	b := store.Get(GetValidatorCurrentRewardsKey(val))
	if b == nil {
		return sdk.NewCoin()
	}
	d.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}


func (d *DistrKeeper) DeleteProposerCurrentRewards(ctx sdk.Context, val sdk.AccAddr) {
	store := ctx.KVStore(d.storeKey)
	store.Delete(GetValidatorCurrentRewardsKey(val))
}

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



func (d *DistrKeeper)GetPreviousProposer(ctx sdk.Context) (proposer sdk.AccAddr){

	store := ctx.KVStore(d.storeKey)
	b := store.Get(types.ProposerKey)
	if b == nil {
		panic("Previous proposer not set")
	}
	d.cdc.MustUnmarshalBinaryLengthPrefixed(b, &proposer)
	return
}


func (d *DistrKeeper)SetPreviousProposer(ctx sdk.Context, proposer sdk.AccAddr) {
	store := ctx.KVStore(d.storeKey)
	b := d.cdc.MustMarshalBinaryLengthPrefixed(proposer)
	store.Set(types.ProposerKey, b)
}

func getKey(val sdk.AccAddr, height int64) sdk.AccAddr {
	add := fmt.Sprintf("%X", val)
	h := strconv.FormatInt(height, 10)
	tKey := add + h
	key := sdk.AccAddr([]byte(tKey))
	return key
}