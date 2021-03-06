package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)
var DistributionCdc *codec.Codec

func init()  {
	DistributionCdc = codec.New()
	RegisterCodec(DistributionCdc)
	codec.RegisterCrypto(DistributionCdc)
	DistributionCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterConcrete(&MsgSetWithdrawAddress{}, "ci123chain/SetWithdrawAddressTx", nil)
	cdc.RegisterConcrete(&MsgWithdrawDelegatorReward{}, "ci123chain/WithdrawDelegatorRewardTx", nil)
	cdc.RegisterConcrete(&MsgWithdrawValidatorCommission{}, "ci123chain/WithdrawValidatorCommissionTx", nil)
	cdc.RegisterConcrete(&MsgFundCommunityPool{}, "ci123chain/FundCommunityPoolTx", nil)
}