package transaction

import "github.com/tendermint/go-amino"

func RegisterCodec(cdc *amino.Codec)  {
	cdc.RegisterInterface((*Transaction)(nil), nil)
	cdc.RegisterConcrete(TransferTx{}, "transaction/transfer", nil)
}