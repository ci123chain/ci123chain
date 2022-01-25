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
	cdc.RegisterConcrete(&MsgRedelegate{}, "ci123chain/MsgPreStakingRedelegate", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "ci123chain/MsgPreStakingUndelegate", nil)
	cdc.RegisterConcrete(&MsgSetStakingToken{}, "ci123chain/MsgSetStakingToken", nil)
	cdc.RegisterConcrete(&VaultRecord{}, "ci123chain/VaultRecord", nil)
	cdc.RegisterConcrete(&Vault{}, "ci123chain/Vault", nil)
	cdc.RegisterConcrete(&StakingRecords{}, "ci123chain/StakingRecords", nil)
	cdc.RegisterConcrete(&MsgPrestakingCreateValidator{}, "ci123chain/MsgPrestakingCreateValidator", nil)
}