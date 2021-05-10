package basic


import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	client "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/infrastructure"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {

}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	return
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator) json.RawMessage {
	return infrastructure.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (am AppModuleBasic) Name() string {
	return infrastructure.ModuleName
}
