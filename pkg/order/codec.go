package order

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

var OrCdc *codec.Codec


func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&UpgradeTx{}, "upgrade/upgrade", nil)
}


func init() {
	OrCdc = codec.New()
	transaction.RegisterCodec(OrCdc)
	RegisterCodec(OrCdc)
}
