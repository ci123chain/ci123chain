package transaction

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/tendermint/go-amino"
)
var ModuleCdc *codec.Codec
func RegisterCodec(cdc *amino.Codec)  {
	cdc.RegisterInterface((*Transaction)(nil), nil)
}

func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}