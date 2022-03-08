package types



import "github.com/ci123chain/ci123chain/pkg/abci/codec"

var PreStakingCodec *codec.Codec

func init(){
	PreStakingCodec = codec.New()
	RegisterCodec(PreStakingCodec)
	codec.RegisterCrypto(PreStakingCodec)
	PreStakingCodec.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	//
	cdc.RegisterConcrete(&MsgPreStaking{}, "ci123chain/MsgPreStaking", nil)
	cdc.RegisterConcrete(&MsgStaking{}, "ci123chain/MsgStaking", nil)
	cdc.RegisterConcrete(&MsgStakingDirect{}, "ci123chain/MsgStakingDirect", nil)

	cdc.RegisterConcrete(&MsgRedelegate{}, "ci123chain/MsgPreStakingRedelegate", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "ci123chain/MsgPreStakingUndelegate", nil)
	cdc.RegisterConcrete(&MsgSetStakingToken{}, "ci123chain/MsgSetStakingToken", nil)
	cdc.RegisterConcrete(&StakingVaultOld{}, "ci123chain/StakingVault", nil)
	cdc.RegisterConcrete(&StakingVault{}, "ci123chain/StakingVaultNew", nil)
	cdc.RegisterConcrete(&TransLog{}, "ci123chain/TransLog", nil)

	cdc.RegisterConcrete(&MsgPrestakingCreateValidator{}, "ci123chain/MsgPrestakingCreateValidator", nil)
	cdc.RegisterConcrete(&MsgPrestakingCreateValidatorDirect{}, "ci123chain/MsgPrestakingCreateValidatorDirect", nil)

}