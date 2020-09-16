package wasm

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
					_, err = wasmer.Instantiate(ctx, code, invoker, 0, args, data.Name, data.Version, data.Author, data.Email, data.Describe, true, address)
					if err != nil {
						panic(err)
					}
				}else if v.Method == types.InvokeMethod {
					_, err = wasmer.Execute(ctx, address, invoker, args)
					if err != nil {
						panic(err)
					}
				}else {
					panic(errors.New(fmt.Sprintf("implement method %s", v.Method)))
				}
			}
		}
	}
}
