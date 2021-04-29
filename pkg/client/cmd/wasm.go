package cmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
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

	wasmCmd.Flags().String(util.FlagAddress, "", "the address of your account")
	wasmCmd.Flags().String(util.FlagGas, "", "expected gas of transaction")
	wasmCmd.Flags().String(util.FlagPrivateKey, "", "the privateKey of account")
	wasmCmd.Flags().String(util.FlagFunds, "", "funds of contract")
	wasmCmd.Flags().String(util.FlagArgs, "", "args of call contract")
	wasmCmd.Flags().String(util.FlagFile, "", "the path of contract file")
	wasmCmd.Flags().String(util.FlagHash, "", "hash of contract code")
	wasmCmd.Flags().String(util.FlagName, "", "name of contract")
	wasmCmd.Flags().String(util.FlagVersion, "", "version of contract")
	wasmCmd.Flags().String(util.FlagAuthor, "", "author of contract")
	wasmCmd.Flags().String(util.FlagEmail, "", "email of contract author")
	wasmCmd.Flags().String(util.FlagDescribe, "", "describe of contract")
	wasmCmd.Flags().String(util.FlagContractAddress, "", "address of contract account")

	util.CheckRequiredFlag(wasmCmd, util.FlagGas)
	util.CheckRequiredFlag(wasmCmd, util.FlagPrivateKey)
	util.CheckRequiredFlag(wasmCmd, util.FlagAddress)
}

var wasmCmd = &cobra.Command{
	Use: "wasm [functionName]",
	Short: "Wasm transaction subcommands",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		funcName := args[0]
		switch funcName {
		case "upload":
			return uploadContract()
		case "init":
			return initContract()
		case "execute":
			return executeContract()
		case "migrate":
			return migrateContract()
		}

		return nil
	},
}

func uploadContract() error {
	ctx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return  err
	}
	from, gas, nonce, key, _, err := GetArgs(ctx)
	if err != nil {
		return err
	}
	fpath := viper.GetString(util.FlagFile)
	fext := path.Ext(fpath)
	if fext != ".wasm" {
		return errors.New("unexpected file")
	}
	code, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	tx, err := sdk.SignUploadContractMsg(code, from, gas, nonce, key)
	if err != nil {
		return err
	}
	txid, err := ctx.BroadcastSignedData(tx)
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
	fpath := viper.GetString(util.FlagFile)
	fext := path.Ext(fpath)
	if fext != ".wasm" {
		return errors.New("unexpected file")
	}
	codeHash := viper.GetString(util.FlagCodeHash)
	name := viper.GetString(util.FlagName)
	version := viper.GetString(util.FlagVersion)
	author := viper.GetString(util.FlagAuthor)
	email := viper.GetString(util.FlagEmail)
	describe := viper.GetString(util.FlagDescribe)

	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		return err
	}
	tx, err := sdk.SignInstantiateContractMsg(hash, from, gas, nonce, key, name, version, author, email, describe, args)
	if err != nil {
		return err
	}
	txid, err := ctx.BroadcastSignedData(tx)
	if err != nil {
		return err
	}
	fmt.Println(txid)
	return nil
}

func executeContract() error {
	ctx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return  err
	}
	from, gas, nonce, key, args, err := GetArgs(ctx)
	if err != nil {
		return err
	}
	contractAddr := viper.GetString(util.FlagContractAddress)
	addrs := types.HexToAddress(contractAddr)
	contractAddress := addrs
	tx, err := sdk.SignExecuteContractMsg(from, gas, nonce, key, contractAddress, args)
	txid, err := ctx.BroadcastSignedData(tx)
	if err != nil {
		return err
	}
	fmt.Println(txid)
	return nil
}

func migrateContract() error {
	ctx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return  err
	}
	from, gas, nonce, key, args, err := GetArgs(ctx)
	if err != nil {
		return err
	}
	fpath := viper.GetString(util.FlagFile)
	fext := path.Ext(fpath)
	if fext != ".wasm" {
		return errors.New("unexpected file")
	}
	codeHash := viper.GetString(util.FlagCodeHash)
	name := viper.GetString(util.FlagName)
	version := viper.GetString(util.FlagVersion)
	author := viper.GetString(util.FlagAuthor)
	email := viper.GetString(util.FlagEmail)
	describe := viper.GetString(util.FlagDescribe)
	contract := viper.GetString(util.FlagContractAddress)
	contractAddr := types.HexToAddress(contract)
	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		return err
	}
	tx, err := sdk.SignMigrateContractMsg(hash, from, gas, nonce, key, name, version, author, email, describe, contractAddr, args)
	if err != nil {
		return err
	}
	txid, err := ctx.BroadcastSignedData(tx)
	if err != nil {
		return err
	}
	fmt.Println(txid)
	return nil
}

func GetArgs(ctx context.Context) (types.AccAddress, uint64, uint64, string, utils.CallData,  error) {
	addrs := viper.GetString(util.FlagAddress)
	address := types.HexToAddress(addrs)

	nonce, _, err := ctx.GetNonceByAddress(address, false)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", utils.CallData{}, err
	}
	gas := viper.GetString(util.FlagGas)
	Gas, err := strconv.ParseUint(gas, 10, 64)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", utils.CallData{}, err
	}
	key := viper.GetString(util.FlagPrivateKey)
	if key == "" {
		return types.AccAddress{}, 0, 0, "", utils.CallData{}, errors.New("privateKey can not be empty")
	}
	Msg := viper.GetString(util.FlagArgs)
	var params utils.CallData
	if Msg == "" {
		return types.AccAddress{}, 0, 0, "", utils.CallData{}, errors.New("calldata can not be empty")
	}else {
		argsByte := []byte(Msg)
		err := json.Unmarshal(argsByte, &params)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", utils.CallData{}, errors.New("unexpected args")
		}
	}
	return address, Gas, nonce, key, params, nil
}