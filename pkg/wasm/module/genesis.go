package module

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	"io/ioutil"
)


func InitGenesis(ctx sdk.Context, wasmer types.WasmKeeperI) {
	var contracts = types.DefaultGenesisState()
	for i := 0; i < len(contracts.Contracts); i++ {
		for _, v := range contracts.Contracts {
			if v.Index == i {
				cdata, _ := base64.StdEncoding.DecodeString(v.Code)
				rdata := bytes.NewReader(cdata)
				r, _ := gzip.NewReader(rdata)
				code, err := ioutil.ReadAll(r)
				if err != nil {
					panic(err)
				}
				address := sdk.HexToAddress(v.Address)
				invoker := sdk.HexToAddress(contracts.Invoker)
				var params types.CallContractParam
				params.Args = v.Params
				params.Method = v.Method
				args, err := json.Marshal(params)
				if err != nil {
					panic(err)
				}
				ctx = ctx.WithValue(types.SystemContract, true)
				if v.Method == types.InitMethod {
					_, err = wasmer.Instantiate(ctx, code, invoker, 0, args, contracts.Name, contracts.Version, contracts.Author, contracts.Email, contracts.Describe, address)
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
