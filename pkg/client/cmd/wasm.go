package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/util"
	wasm "github.com/tanhuiya/ci123chain/pkg/wasm/types"
	sdk "github.com/tanhuiya/ci123chain/sdk/wasm"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

func init() {
	rootCmd.AddCommand(WasmCmd)

	WasmCmd.Flags().String(helper.FlagAddress, "", "the address of your account")
	WasmCmd.Flags().String(helper.FlagGas, "", "expected gas of transaction")
	WasmCmd.Flags().String(helper.FlagPrivateKey, "", "the privateKey of account")
	WasmCmd.Flags().String(helper.FlagFunds, "", "funds of contract")
	WasmCmd.Flags().String(helper.FlagArgs, "", "args of call contract")
	WasmCmd.Flags().String(helper.FlagFile, "", "the path of contract file")
	WasmCmd.Flags().String(helper.FlagID, "", "id of contract code")
	WasmCmd.Flags().String(helper.FlagLabel, "", "label of contract")
	WasmCmd.Flags().String(helper.FlagContractAddress, "", "address of contract account")

	util.CheckRequiredFlag(WasmCmd, helper.FlagGas)
	util.CheckRequiredFlag(WasmCmd, helper.FlagPrivateKey)
	util.CheckRequiredFlag(WasmCmd, helper.FlagAddress)
	err := viper.BindPFlags(WasmCmd.Flags())
	if err != nil {
		panic(err)
	}
}

var WasmCmd = &cobra.Command{
	Use: "wasm [functionName]",
	Short: "Wasm transaction subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {

		funcName := args[0]
		switch funcName {
		case "install":
			return installContract()
		case "init":
			return initContract()
		case "invoke":
			return invokeContract()
		}

		return nil
	},
}

func installContract() error {
	ctx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return  err
	}
	fpath := viper.GetString(helper.FlagFile)
	fext := path.Ext(fpath)
	if fext != ".wasm" {
		return errors.New("unexpected file")
	}
	code, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	if ok := wasm.IsValidaWasmFile(code); ok != nil {
		return ok
	}
	from, gas, nonce, key, _, err := GetArgs(ctx)
	if err != nil {
		return err
	}
	txByte, err := sdk.SignStoreCodeMsg(from, gas, nonce, key, from, code)
	txid, err := ctx.BroadcastSignedData(txByte)
	if err != nil {
		return err
	}
	fmt.Println(txid)
	return nil
}

func initContract() error {
	ctx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return  err
	}
	from, gas, nonce, key, args, err := GetArgs(ctx)
	if err != nil {
		return err
	}
	id := viper.GetString(helper.FlagID)
	codeID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	label := viper.GetString(helper.FlagLabel)
	if label == "" {
		label = "demo contract"
	}

	txByte, err := sdk.SignInstantiateContractMsg(from, gas, nonce, codeID, key, from, label, args)
	txid, err := ctx.BroadcastSignedData(txByte)
	if err != nil {
		return err
	}
	fmt.Println(txid)
	return nil
}

func invokeContract() error {
	ctx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return  err
	}
	from, gas, nonce, key, args, err := GetArgs(ctx)
	if err != nil {
		return err
	}
	contractAddr := viper.GetString(helper.FlagContractAddress)
	addrs := types.HexToAddress(contractAddr)
	contractAddress := addrs
	txByte, err := sdk.SignExecuteContractMsg(from, gas, nonce, key, from, contractAddress, args)
	txid, err := ctx.BroadcastSignedData(txByte)
	if err != nil {
		return err
	}
	fmt.Println(txid)
	return nil
}


func GetArgs(ctx context.Context) (types.AccAddress, uint64, uint64, string, json.RawMessage,  error) {
	var args json.RawMessage
	addrs := viper.GetString(helper.FlagAddress)
	address := types.HexToAddress(addrs)

	nonce, err := ctx.GetNonceByAddress(address)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", nil, err
	}
	gas := viper.GetString(helper.FlagGas)
	Gas, err := strconv.ParseUint(gas, 10, 64)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", nil, err
	}
	key := viper.GetString(helper.FlagPrivateKey)
	if key == "" {
		return types.AccAddress{}, 0, 0, "", nil, errors.New("privateKey can not be empty")
	}
	Msg := viper.GetString(helper.FlagArgs)
	if Msg == "" {
		args = json.RawMessage{}
	}else {
		msgStr := getStr(Msg)
		args, err = json.Marshal(msgStr)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, err
		}
	}
	return address, Gas, nonce, key, args, nil
}

func getStr(m string) wasm.CallContractParam {
	var str []string
	var args []string
	var p wasm.CallContractParam
	b := strings.Split(m, "{")
	if b[1] == "" || len(b) != 2 {
		p = wasm.CallContractParam{
			Method: "",
			Args:   nil,
		}
		return p
	}else {
		c := strings.Split(b[1], "}")
		if c[0] == "" || len(c) != 2 {
			p = wasm.CallContractParam{
				Method: "",
				Args:   nil,
			}
			return p
		}else {
			d := strings.Split(c[0], ",")
			if len(d) == 0 {
				p = wasm.CallContractParam{
					Method: "",
					Args:   nil,
				}
				return p
			}
			for i := 0; i < len(d); i++ {
				e := strings.Split(d[i], "\"")
				if len(e) != 3 {
					str = append(str, "")
				}else {
					str = append(str, e[1])
				}
			}
		}
	}
	if len(str) == 0 {
		p = wasm.CallContractParam{
			Method: "",
			Args:   nil,
		}
	}else {
		method := str[0]
		for i := 1; i < len(str); i++ {
			args = append(args, str[i])
		}
		p = wasm.CallContractParam{
			Method: method,
			Args:   args,
		}
	}
	return p
}