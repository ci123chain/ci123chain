package types

import (
	"bytes"
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/app/module"
	infratypes "github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/upgrade"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

var cdc *codec.Codec

// attempt to make some pretty json
func MarshalJSONIndent(cdc *amino.Codec, obj interface{}) ([]byte, error) {
	bz, err := cdc.MarshalJSON(obj)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}

	var out bytes.Buffer
	err = json.Indent(&out, bz, "", "  ")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return out.Bytes(), nil
}

func GetCodec() *codec.Codec {
	if cdc == nil {
		cdc = amino.NewCodec()
		cdc.RegisterConcrete(&CommonTx{}, "transfer/commontx", nil)
		cdc.RegisterConcrete(&MsgEthereumTx{}, "eth/msgEthereumTx", nil)
		sdk.RegisterCodec(cdc)
		transaction.RegisterCodec(cdc)
		transfer.RegisterCodec(cdc)
		module.ModuleBasics.RegisterCodec(cdc)
		infratypes.RegisterCodec(cdc)
		upgrade.RegisterCodec(cdc)
		cryptoAmino.RegisterAmino(cdc)
	}
	return cdc
}

var encodingConfig *EncodingConfig

type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Marshaler
	//TxConfig          client.TxConfig
	Amino             *codec.Codec
}

func GetEncodingConfig() *EncodingConfig{
	if encodingConfig == nil {
		cdc := GetCodec()
		interfaceRegistry :=  types.NewInterfaceRegistry()
		marshaler := codec.NewProtoCodec(interfaceRegistry)

		encodingConfig = &EncodingConfig{
			InterfaceRegistry: interfaceRegistry,
			Marshaler:         marshaler,
			//TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
			Amino:             cdc,
		}
		sdk.RegisterInterfaces(encodingConfig.InterfaceRegistry)
		module.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	}

	return encodingConfig
}
