package cmd

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/client/types"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"github.com/tanhuiya/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo 		= "to"
	flagAmount  = "amount"
	flagGas 	= "gas"
)

func init()  {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().String(flagTo, "", "Address sending to")
	transferCmd.Flags().Uint(flagAmount, 0, "Amount tbe spent")
	transferCmd.Flags().Uint(flagGas, 0, "gas for tx")
	transferCmd.Flags().String(helper.FlagAddress, "", "Address to sign with")
	transferCmd.Flags().String(flagPassword, "", "passphrase")
	util.CheckRequiredFlag(transferCmd, flagAmount)
	util.CheckRequiredFlag(transferCmd, flagGas)
}

var transferCmd = &cobra.Command{
	Use: "transfer",
	Short: "Build, Sign, and send transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return err
		}
		addrs, err := ctx.GetInputAddresses()
		if err != nil {
			return nil
		}
		from := addrs[0]
		tos, err := helper.ParseAddrs(viper.GetString(flagTo))
		if err != nil {
			return types.ErrParseAddr(types.DefaultCodespace, err)
		}
		if len(tos) == 0 {
			return types.ErrNoAddr(types.DefaultCodespace, err)
		}
		//直接getNonce
		//todo err
		nonce, err := ctx.GetNonceByAddress(from)
		if err != nil {
			return err
		}

		ucoin := uint64(viper.GetInt(flagAmount))
		tx := transfer.NewTransferTx(from, tos[0], uint64(viper.GetInt(flagGas)), nonce, sdk.Coin(ucoin), false)

		password := viper.GetString(flagPassword)
		if len(password) < 1 {
			var err error
			password, err = helper.GetPassphrase(from)
			if err != nil {
				return types.ErrGetPassPhrase(types.DefaultCodespace, err)
			}
		}

		signedData, err := getSignedDataWithTx(ctx, tx, password, from)
		if err != nil {
			return types.ErrGetSignData(types.DefaultCodespace, err)
		}
		res, err := ctx.BroadcastSignedData(signedData)
		if err != nil {
			return types.ErrBroadcast(types.DefaultCodespace, err)
		}
		ctx.PrintOutput(res)
		return nil
	},
}

