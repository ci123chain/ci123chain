package cmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/util"
	wasm "github.com/ci123chain/ci123chain/pkg/wasm/types"
	sdk "github.com/ci123chain/ci123chain/sdk/wasm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

func init() {
	rootCmd.AddCommand(wasmCmd)

	wasmCmd.Flags().String(helper.FlagAddress, "", "the address of your account")
	wasmCmd.Flags().String(helper.FlagGas, "", "expected gas of transaction")
	wasmCmd.Flags().String(helper.FlagPrivateKey, "", "the privateKey of account")
	wasmCmd.Flags().String(helper.FlagFunds, "", "funds of contract")
	wasmCmd.Flags().String(helper.FlagArgs, "", "args of call contract")
	wasmCmd.Flags().String(helper.FlagFile, "", "the path of contract file")
	wasmCmd.Flags().String(helper.FlagHash, "", "hash of contract code")
	wasmCmd.Flags().String(helper.FlagLabel, "", "label of contract")
	wasmCmd.Flags().String(helper.FlagContractAddress, "", "address of contract account")

	util.CheckRequiredFlag(wasmCmd, helper.FlagGas)
	util.CheckRequiredFlag(wasmCmd, helper.FlagPrivateKey)
	util.CheckRequiredFlag(wasmCmd, helper.FlagAddress)
}

var wasmCmd = &cobra.Command{
	Use: "wasm [functionName]",
	Short: "Wasm transaction subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

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
	hash := viper.GetString(helper.FlagHash)
	Hash, err := hex.DecodeString(strings.ToLower(hash))
	if err != nil {
		return errors.New("decode codeHash fail")
	}
	label := viper.GetString(helper.FlagLabel)
	if label == "" {
		label = "demo contract"
	}

	txByte, err := sdk.SignInstantiateContractMsg(from, gas, nonce, Hash, key, from, label, args)
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

	nonce, _, err := ctx.GetNonceByAddress(address)
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
		var params wasm.CallContractParam
		argsByte := []byte(Msg)
		err := json.Unmarshal(argsByte, params)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, errors.New("unexpected args")
		}
		args = argsByte
	}
	return address, Gas, nonce, key, args, nil
}