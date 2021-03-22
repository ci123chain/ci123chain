package types

import (
	"bytes"
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app/module"
	ibc "github.com/ci123chain/ci123chain/pkg/ibc/core/types"
	i_types "github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	"github.com/ci123chain/ci123chain/pkg/mortgage"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// attempt to make some pretty json
func MarshalJSONIndent(cdc *amino.Codec, obj interface{}) ([]byte, error) {
	bz, err := cdc.MarshalJSON(obj)
	if err != nil {
		return nil, abci.ErrInternal("Marshal failed")
	}

	var out bytes.Buffer
	err = json.Indent(&out, bz, "", "  ")
	if err != nil {
		return nil, abci.ErrInternal("Indent failed")
	}
	return out.Bytes(), nil
}

func MakeCodec() *codec.Codec {
	cdc := amino.NewCodec()
	//cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	//cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(&CommonTx{}, "transfer/commontx", nil)
	cdc.RegisterConcrete(&MsgEthereumTx{}, "eth/msgEthereumTx", nil)
	sdk.RegisterCodec(cdc)
	transaction.RegisterCodec(cdc)
	transfer.RegisterCodec(cdc)
	mortgage.RegisterCodec(cdc)
	module.ModuleBasics.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)
	i_types.RegisterCodec(cdc)
	//acc_types.RegisterCodec(cdc)
	//cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
	//	secp256k1.PubKeyAminoName, nil)
	//cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{},
	//	secp256k1.PrivKeyAminoName, nil)
	cryptoAmino.RegisterAmino(cdc)
	return cdc
}
