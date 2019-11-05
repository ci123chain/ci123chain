package account

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account/keeper"
)


type AppModule struct {
	AccountKeeper	keeper.AccountKeeper
}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage)  {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, ModuleCdc, am.AccountKeeper, genesisState)
}
