package db

import "github.com/tendermint/go-amino"

var cdc = amino.NewCodec()

func init() {
	RegisterAmino(cdc)
}

func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(RWSet{},
		"ci123chain/RWSet", nil)
	cdc.RegisterConcrete(RWSetItems{},
		"ci123chain/RWSetItem", nil)
}

