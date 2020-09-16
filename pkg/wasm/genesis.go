package wasm

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/wasm/keeper"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
)

const (
	gasWanted = uint64(80000000)
)

func InitGenesis(ctx sdk.Context, wasmer keeper.Keeper, data GenesisState) {

	for i := 0; i < len(data.Contracts); i++ {
		for _, v := range data.Contracts {
			if v.Index == i {
				code, err := hex.DecodeString(v.Code)
				if err != nil {
					panic(err)
				}
				address := sdk.HexToAddress(v.Address)
				invoker := sdk.HexToAddress(data.Invoker)
				var params types.CallContractParam
				params.Args = v.Params
				params.Method = v.Method
				args, err := json.Marshal(params)
				if err != nil {
					panic(err)
				}
				keeper.SetGasWanted(gasWanted)
				if v.Method == types.InitMethod {
					err = wasmer.GenesisContractInit(ctx, code, invoker,args, data.Name, data.Version, data.Author, data.Email, data.Describe, address)
					if err != nil {
						panic(err)
					}
				}else {
					err = wasmer.GenesisInvoke(ctx, address, invoker, args)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
