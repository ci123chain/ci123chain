package cmd

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/client/types"
	transfer2 "github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagFrom    = "from"
	flagTo 		= "to"
	flagAmount  = "amount"
	flagGas 	= "gas"
	flagKey		= "privKey"
	flagIsFabric= "isFabric"
	flagDenom   = "denom"
)

func init()  {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().String(flagTo, "", "Address sending to")
	transferCmd.Flags().Uint(flagAmount, 0, "Amount tbe spent")
	transferCmd.Flags().Uint(flagGas, 0, "gas for tx")
	transferCmd.Flags().String(helper.FlagAddress, "", "Address to sign with")
	transferCmd.Flags().String(flagPassword, "", "passphrase")
	transferCmd.Flags().String(flagDenom, "", "coin denom")

	util.CheckRequiredFlag(transferCmd, flagAmount)
	util.CheckRequiredFlag(transferCmd, flagGas)
}

var transferCmd = &cobra.Command{
	Use: "transfer",
	Short: "Build, Sign, and send transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			panic(err)
		}

		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return err
		}
		from := sdk.HexToAddress(viper.GetString(flagFrom))
		tos, err := helper.ParseAddrs(viper.GetString(flagTo))
		if err != nil {
			return types.ErrParseAddr(types.DefaultCodespace, err)
		}
		if len(tos) == 0 {
			return types.ErrNoAddr(types.DefaultCodespace, err)
		}
		d := viper.GetString(flagDenom)
		if d == "" {
			return types.ErrParseParam(types.DefaultCodespace, errors.New("coin denom can't be emtpy"))
		}

		gas := uint64((viper.GetInt(flagGas)))
		amount := uint64(viper.GetInt(flagAmount))
		privKey := viper.GetString(flagKey)
		isFabric := viper.GetBool(flagIsFabric)

		coin := sdk.NewUInt64Coin(d, amount)
		msg := transfer2.NewMsgTransfer(from, tos[0], coin, isFabric)
		nonce, err := transfer2.GetNonceByAddress(from)
		if err != nil {
			return types.ErrParseParam(types.DefaultCodespace, err)
		}

		txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
		if err != nil {
			return types.ErrParseParam(types.DefaultCodespace, err)
		}

		res, err := ctx.BroadcastSignedData(txByte)
		if err != nil {
			return types.ErrBroadcast(types.DefaultCodespace, err)
		}
		ctx.PrintOutput(res)
		return nil
	},
}

