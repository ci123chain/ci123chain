package types


import "github.com/ci123chain/ci123chain/pkg/abci/codec"

var InfrastructureCdc *codec.Codec


func init() {
	InfrastructureCdc = codec.New()
	RegisterCodec(InfrastructureCdc)
	codec.RegisterCrypto(InfrastructureCdc)
	InfrastructureCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgStoreContent{}, "ci123chain/StoreContentTx", nil)
}
