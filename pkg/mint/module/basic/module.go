package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	client "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/mint"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModuleBasic struct {}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	mint.RegisterCodec(codec)
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	return
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator) json.RawMessage {
	return mint.ModuleCdc.MustMarshalJSON(mint.DefaultGenesisState())
}

func (am AppModuleBasic) Name() string {
	return mint.ModuleName
}

