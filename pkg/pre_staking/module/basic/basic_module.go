package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModuleBasic struct {}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(context.Context, *runtime.ServeMux) {
}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	return
}




func (am AppModuleBasic) DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage {
	p := types.DefaultGenesisState(validators)
	b, err := types.PreStakingCodec.MarshalJSONIndent(p, "", "")
	if err != nil{
		panic(err)
	}
	return b
}


func (am AppModuleBasic) Name() string {
	return types.ModuleName
}
