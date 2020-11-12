package module

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"io/ioutil"
)


func WasmInitGenesis(ctx sdk.Context, wasmer moduletypes.KeeperI) {
	var contracts = types.DefaultGenesisState()
	ctx = ctx.WithValue(types.SystemContract, true)
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
				if v.Method == types.InitMethod {
					codeHash, err := wasmer.Upload(ctx, code, invoker)
					if err != nil {
						panic(err)
					}
					_, err = wasmer.Instantiate(ctx, codeHash, invoker, args, contracts.Name, contracts.Version, contracts.Author, contracts.Email, contracts.Describe, address)

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

func EvmInitGenesis(ctx sdk.Context, k moduletypes.KeeperI, data evmtypes.GenesisState) {
	for _, account := range data.Accounts {
		// FIXME: this will override bank InitGenesis balance!
		k.SetBalance(ctx, account.Address, account.Balance)
		k.SetCode(ctx, account.Address, account.Code)
		for _, storage := range account.Storage {
			k.SetState(ctx, account.Address, storage.Key, storage.Value)
		}
	}

	var err error
	for _, txLog := range data.TxsLogs {
		err = k.SetLogs(ctx, txLog.Hash, txLog.Logs)
		if err != nil {
			panic(err)
		}
	}

	k.SetChainConfig(ctx, data.ChainConfig)
	k.SetParams(ctx, data.Params)

	// set state objects and code to store
	_, err = k.Commit(ctx, false)
	if err != nil {
		panic(err)
	}

	// set storage to store
	// NOTE: don't delete empty object to prevent import-export simulation failure
	err = k.Finalise(ctx, false)
	if err != nil {
		panic(err)
	}

	return
}


