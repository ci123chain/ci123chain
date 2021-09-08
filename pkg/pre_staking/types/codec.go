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
	//cdc.RegisterConcrete(&MsgDelegate{}, "ci123chain/MsgDelegate", nil)
	//cdc.RegisterConcrete(&MsgUndelegate{}, "ci123chain/MsgUndelegate", nil)
	//cdc.RegisterConcrete(&MsgRedelegate{}, "ci123chain/MsgRedelegate", nil)
}