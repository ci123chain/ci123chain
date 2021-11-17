package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	client "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/upgrade/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	return
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

func (am AppModuleBasic) DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage {
	return json.RawMessage{}
}

func (am AppModuleBasic) Name() string {
	return types.ModuleName
}