package module

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/order"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModuleGenesis interface {
	AppModuleBasic

	// 根据 genesis 配置初始化
	InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate
	BeginBlocker(ctx types.Context, req abci.RequestBeginBlock)
	Committer(ctx types.Context)
	EndBlock(ctx types.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate
}

type AppModuleBasic interface {
	Name() string
	RegisterCodec(codec *codec.Codec)

	// 默认的 genesis 配置
	DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage
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

func (bm BasicManager) DefaultGenesis(validators []tmtypes.GenesisValidator) map[string]json.RawMessage {
	genesis := make(map[string]json.RawMessage)
	for _, b := range bm {
		if b.DefaultGenesis(validators) != nil {
			genesis[b.Name()] = b.DefaultGenesis(validators)
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
	o, ok := am.Modules[order.ModuleName]
	if ok && o != nil{
		o.BeginBlocker(ctx, req)
	}

	for _, m := range am.Modules {
		if m == am.Modules[order.ModuleName] {
			continue
		}
		m.BeginBlocker(ctx, req)
	}
	return abci.ResponseBeginBlock{}
}

func (am AppManager) Committer(ctx types.Context) abci.ResponseCommit {

	for _, m := range am.Modules {
		m.Committer(ctx)
	}
	return abci.ResponseCommit{}
}

func (am AppManager) EndBlocker(ctx types.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {

	validatorUpdates := []abci.ValidatorUpdate{}
	for _, m := range am.Modules {
		if m == am.Modules[order.ModuleName] {
			continue
		}
		moduleValUpdates := m.EndBlock(ctx, req)
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator EndBlock updates already set by a previous module")
			}

			validatorUpdates = moduleValUpdates
		}
	}
	return abci.ResponseEndBlock{ValidatorUpdates:validatorUpdates}
}