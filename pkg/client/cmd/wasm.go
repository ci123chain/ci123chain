package cmd

import (
	"encoding/hex"
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
	sdk "github.com/tanhuiya/ci123chain/sdk/wasm"
	"io/ioutil"
	"strconv"
)

func init() {
	WasmCmd.AddCommand(StoreCodeCmd,
		InstantiateContractCmd,
		ExecuteContractCmd)
	rootCmd.AddCommand(WasmCmd)

	StoreCodeCmd.Flags().String(helper.FlagFile, "", "the path of contract file")
	StoreCodeCmd.Flags().String(helper.FlagGas, "", "expected gas of transaction")
	StoreCodeCmd.Flags().String(helper.FlagPrivateKey, "", "the privateKey of account")

	InstantiateContractCmd.Flags().String(helper.FlagGas, "", "expected gas of transaction")
	InstantiateContractCmd.Flags().String(helper.FlagPrivateKey, "", "the privateKey of account")
	InstantiateContractCmd.Flags().String(helper.FlagID, "", "id of contract code")
	InstantiateContractCmd.Flags().String(helper.FlagLabel, "", "label of contract")
	InstantiateContractCmd.Flags().String(helper.FlagFunds, "", "funds of contract")
	InstantiateContractCmd.Flags().String(helper.FlagMsg, "", "message of init contract")

	ExecuteContractCmd.Flags().String(helper.FlagGas, "", "expected gas of transaction")
	ExecuteContractCmd.Flags().String(helper.FlagPrivateKey, "", "the privateKey of account")
	ExecuteContractCmd.Flags().String(helper.FlagContractAddress, "", "address of contract account")
	ExecuteContractCmd.Flags().String(helper.FlagMsg, "", "msg of execute contract")
	ExecuteContractCmd.Flags().String(helper.FlagFunds, "", "funds of contract")

	util.CheckRequiredFlag(StoreCodeCmd, helper.FlagFile)
	util.CheckRequiredFlag(StoreCodeCmd, helper.FlagGas)
	util.CheckRequiredFlag(StoreCodeCmd, helper.FlagPrivateKey)

	util.CheckRequiredFlag(InstantiateContractCmd, helper.FlagGas)
	util.CheckRequiredFlag(InstantiateContractCmd, helper.FlagPrivateKey)
	util.CheckRequiredFlag(InstantiateContractCmd, helper.FlagID)
	util.CheckRequiredFlag(InstantiateContractCmd, helper.FlagLabel)
	util.CheckRequiredFlag(InstantiateContractCmd, helper.FlagFunds)
	util.CheckRequiredFlag(InstantiateContractCmd, helper.FlagMsg)

	util.CheckRequiredFlag(ExecuteContractCmd, helper.FlagGas)
	util.CheckRequiredFlag(ExecuteContractCmd, helper.FlagPrivateKey)
	util.CheckRequiredFlag(ExecuteContractCmd, helper.FlagContractAddress)
	util.CheckRequiredFlag(ExecuteContractCmd, helper.FlagMsg)
	util.CheckRequiredFlag(ExecuteContractCmd, helper.FlagFunds)
}


var StoreCodeCmd = &cobra.Command{
	Use: "store",
	Short: "Upload a wasm binary",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return  err
		}
		code, err := ioutil.ReadFile(helper.FlagFile)
		if err != nil {
			return err
		}
		from, gas, nonce, key, _, _, err := GetArgs(ctx)
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
	},
}

var InstantiateContractCmd = &cobra.Command{
	Use: "init",
	Short: "Init a wasm contract",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return  err
		}
		from, gas, nonce, key, funds, msg, err := GetArgs(ctx)
		if err != nil {
			return err
		}
		id := viper.GetString(helper.FlagID)
		codeID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return err
		}
		label := viper.GetString(helper.FlagLabel)

		txByte, err := sdk.SignInstantiateContractMsg(from, gas, nonce, codeID, key, from, label, msg, funds)
		txid, err := ctx.BroadcastSignedData(txByte)
		if err != nil {
			return err
		}
		fmt.Println(txid)

		return nil
	},
}

var ExecuteContractCmd = &cobra.Command{
	Use: "invoke",
	Short: "Execute a command on a wasm contract",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return  err
		}
		from, gas, nonce, key, funds, msg, err := GetArgs(ctx)
		if err != nil {
			return err
		}
		contractAddr := viper.GetString(helper.FlagAddress)
		addrs, err := helper.ParseAddrs(contractAddr)
		if err != nil {
			return err
		}
		contractAddress := addrs[0]
		txByte, err := sdk.SignExecuteContractMsg(from, gas, nonce, key, from, contractAddress, msg, funds)
		txid, err := ctx.BroadcastSignedData(txByte)
		if err != nil {
			return err
		}
		fmt.Println(txid)
		return nil
	},
}

var WasmCmd = &cobra.Command{
	Use: "wasm",
	Short: "Wasm transaction subcommands",
}


func GetArgs(ctx context.Context) (types.AccAddress, uint64, uint64, string, types.Coin, json.RawMessage,  error) {
	var JsonMsg interface{}
	addrs, err := ctx.GetInputAddresses()
	if err != nil {
		return types.AccAddress{}, 0, 0, "", types.Coin{}, nil, err
	}
	nonce, err := ctx.GetNonceByAddress(addrs[0])
	if err != nil {
		return types.AccAddress{}, 0, 0, "", types.Coin{}, nil, err
	}
	gas := viper.GetString(helper.FlagGas)
	Gas, err := strconv.ParseUint(gas, 10, 64)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", types.Coin{}, nil, err
	}
	key := viper.GetString(helper.FlagPrivateKey)
	if key == "" {
		return types.AccAddress{}, 0, 0, "", types.Coin{}, nil, errors.New("privateKey can not be empty")
	}

	funds := viper.GetString(helper.FlagFunds)
	fs, err := strconv.ParseInt(funds, 10, 64)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", types.Coin{}, nil, errors.New("privateKey can not be empty")
	}
	Funds := types.NewCoin(types.NewInt(fs))
	Msg := viper.GetString(helper.FlagMsg)
	msgByte, err := hex.DecodeString(Msg)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", types.Coin{}, nil, err
	}
	msg := json.RawMessage(msgByte)
	err = json.Unmarshal(msgByte, &JsonMsg)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", types.Coin{}, nil, err
	}
	return addrs[0], Gas, nonce, key, Funds, msg, nil
}