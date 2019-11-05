package app

import (
	"bytes"
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// attempt to make some pretty json
func MarshalJSONIndent(cdc *amino.Codec, obj interface{}) ([]byte, error) {
	bz, err := cdc.MarshalJSON(obj)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = json.Indent(&out, bz, "", "  ")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func  MakeCodec() *codec.Codec {
	cdc := amino.NewCodec()
	//cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	//cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	sdk.RegisterCodec(cdc)
	transaction.RegisterCodec(cdc)

	ModuleBasics.RegisterCodec(cdc)
	//acc_types.RegisterCodec(cdc)
	//cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
	//	secp256k1.PubKeyAminoName, nil)
	//cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{},
	//	secp256k1.PrivKeyAminoName, nil)
	cryptoAmino.RegisterAmino(cdc)
	return cdc
}
