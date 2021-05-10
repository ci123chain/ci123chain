package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	client "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {}


func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
	evmtypes.RegisterCodec(codec)
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	return
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

func (am AppModuleBasic) Name() string {
	return moduletypes.ModuleName
}

// DefaultGenesis is json default structure
func (am AppModuleBasic) DefaultGenesis(vals []tmtypes.GenesisValidator) json.RawMessage {
	return evmtypes.ModuleCdc.MustMarshalJSON(evmtypes.DefaultGenesisState())
}

// ValidateGenesis is the validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var genesisState evmtypes.GenesisState
	err := evmtypes.ModuleCdc.UnmarshalJSON(bz, &genesisState)
	if err != nil {
		return err
	}

	return genesisState.Validate()
}