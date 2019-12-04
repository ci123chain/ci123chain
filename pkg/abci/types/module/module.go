package module

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModuleGenesis interface {
	AppModuleBasic

	// 根据 genesis 配置初始化
	InitGenesis(ctx types.Context, data json.RawMessage)
	BeginBlocker(ctx types.Context, req abci.RequestBeginBlock)
}

type AppModuleBasic interface {
	Name() string
	RegisterCodec(codec *codec.Codec)

	// 默认的 genesis 配置
	DefaultGenesis() json.RawMessage
}


// BasicManager
type BasicManager map[string]AppModuleBasic

func NewBasicManager(modules ...AppModuleBasic) BasicManager {
	moduleMap := make(map[string]AppModuleBasic)
	for _, module := range modules {
		moduleMap[module.Name()] = module
	}
	return moduleMap
}

func (bm BasicManager)RegisterCodec(cdc *codec.Codec)  {
	for _, b := range bm {
		b.RegisterCodec(cdc)
	}
}

func (bm BasicManager) DefaultGenesis() map[string]json.RawMessage {
	genesis := make(map[string]json.RawMessage)
	for _, b := range bm {
		if b.DefaultGenesis() != nil {
			genesis[b.Name()] = b.DefaultGenesis()
		}
	}
	return genesis
}



type AppModule interface {
	AppModuleGenesis
}

type AppManager struct {
	Modules 	map[string]AppModule
}

func NewManager(modules ...AppModule) *AppManager {
	moduleMap := make(map[string]AppModule)
	for _, module := range modules {
		moduleMap[module.Name()] = module
	}
	return &AppManager{
		Modules: moduleMap,
	}
}

func (am AppManager) InitGenesis(ctx types.Context, data map[string]json.RawMessage) abci.ResponseInitChain {
	var validatorUpdates []abci.ValidatorUpdate
	for _, m := range am.Modules {
		m.InitGenesis(ctx, data[m.Name()])
	}
	return abci.ResponseInitChain{
		Validators: validatorUpdates,
	}
}

func (am AppManager) BeginBlocker(ctx types.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {

	for _, m := range am.Modules {
		m.BeginBlocker(ctx, req)
	}
	return abci.ResponseBeginBlock{}
}