package distribution
//
//import (
//	"github.com/ci123chain/ci123chain/pkg/abci/codec"
//	"github.com/ci123chain/ci123chain/pkg/abci/types"
//	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
//	"github.com/ci123chain/ci123chain/pkg/account"
//)
//
//// keeper of the staking store
//type DistrKeeper struct {
//	storeKey            types.StoreKey
//	cdc                 *codec.Codec
//	feeCollectionKeeper fc.FcKeeper
//	ak                  account.AccountKeeper
//}
//
//var (
//	ValidatorCurrentRewardsPrefix        = []byte{0x06}
//	DisrtKey = "distr"
//)
//
//// create a new keeper
//func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, fck fc.FcKeeper, ak account.AccountKeeper) DistrKeeper {
//	keeper := DistrKeeper{
//		storeKey:            key,
//		cdc:                 cdc,
//		feeCollectionKeeper: fck,
//		ak:                  ak,
//	}
//	return keeper
//}
//
//func GetValidatorCurrentRewardsKey(v types.AccAddress) []byte {
//	return append(ValidatorCurrentRewardsPrefix, v.Bytes()...)
//}
//
//func (d *DistrKeeper) SetProposerCurrentRewards(ctx types.Context, val types.AccAddress, rewards types.Coin) {
//	store := ctx.KVStore(d.storeKey)
//	b := d.cdc.MustMarshalBinaryLengthPrefixed(rewards)
//	store.Set(GetValidatorCurrentRewardsKey(val), b)
//}
//
//func (d *DistrKeeper) GetProposerCurrentRewards(ctx types.Context, val types.AccAddress) (rewards types.Coin) {
//	store := ctx.KVStore(d.storeKey)
//	b := store.Get(GetValidatorCurrentRewardsKey(val))
//	if b == nil {
//		return types.NewCoin()
//	}
//	d.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
//	return
//}
//
//func (d *DistrKeeper) DeleteProposerCurrentRewards(ctx types.Context, val types.AccAddress) {
//	store := ctx.KVStore(d.storeKey)
//	store.Delete(GetValidatorCurrentRewardsKey(val))
//}
//
//func (d *DistrKeeper) DistributeRewardsToValidators(ctx types.Context, proposer types.AccAddress, fee types.Coin) {
//
//	account := d.ak.GetAccount(ctx, proposer)
//	accCoin := account.GetCoin()
//	accCoin.SafeAdd(fee)
//	d.ak.SetAccount(ctx, account)
//	//n := float32(0.05)
//	//mulNum := length + n
//	//var v float32
//	//var mulFee = float32(fee)
//	//v = mulFee/mulNum
//	//fmt.Print(v)
//	//val := v * 0.05
//	//value := types.Coin(uint64(val))
//	//proposerAcc := d.ak.GetAccount(ctx, proposer)
//	//accCoin := proposerAcc.GetCoin()
//	//accCoin.SafeAdd(value)
//	//d.ak.SetAccount(ctx, proposerAcc)
//	//
//	//validatorVal := types.Coin(uint64(v))
//	//for i, _ := range validators {
//	//	validatorAcc := d.ak.GetAccount(ctx, validators[i])
//	//	accCoin := validatorAcc.GetCoin()
//	//	accCoin.SafeAdd(validatorVal)
//	//	d.ak.SetAccount(ctx, validatorAcc)
//	//}
//}
//
