package transaction

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tendermint/go-amino"
)

func RegisterCodec(cdc *amino.Codec)  {
	cdc.RegisterInterface((*Transaction)(nil), nil)
	cdc.RegisterConcrete(&CommonTx{}, "transaction/commontx", nil)
	cdc.RegisterConcrete(&TransferTx{}, "transaction/transfer", nil)
}

var transferCdc *codec.Codec

func init()  {
	transferCdc = codec.New()
	transferCdc.RegisterInterface((*types.Tx)(nil), nil)
	RegisterCodec(transferCdc)
}