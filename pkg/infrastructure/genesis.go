package infrastructure

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/keeper"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
)


func InitGenesis(ctx sdk.Context, k keeper.InfrastructureKeeper, data types.GenesisState) {
	for _, v := range data.Data {
		value, _ := json.Marshal(v)
		k.SetContent(ctx, []byte(v.Key), value)
	}
}



func ExportGenesis(ctx sdk.Context, k keeper.InfrastructureKeeper) types.GenesisState {

	store := ctx.KVStore(k.GetStoreKey())
	iter := sdk.KVStorePrefixIterator(store, nil)

	defer iter.Close()
	var data = make([]types.StoredContent, 0)
	for ; iter.Valid(); iter.Next() {
		var d types.StoredContent
		v := iter.Value()
		err := json.Unmarshal(v, &d)
		if err != nil {
			panic(err)
		}
		data = append(data, d)
	}

	return types.GenesisState{
		Data:data,
	}
}