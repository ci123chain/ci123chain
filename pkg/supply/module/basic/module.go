package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	client "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	types2 "github.com/ci123chain/ci123chain/pkg/supply/types"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {
}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	supply.RegisterCodec(codec)
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	return
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator) json.RawMessage {
	var res = types2.DefaultGenesisState()
	return supply.ModuleCdc.MustMarshalJSON(res)

}

func (am AppModuleBasic) Name() string {
	return supply.ModuleName
}