package cmd

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
	"errors"
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
	Short: "Build, Sign, and send transaction",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx, err := client.NewClientContextFromViper()
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
			return err
		}
		if len(tos) == 0 {
			return errors.New("must provide an address to send to")
		}
		nonce, err := transaction.GetNonceByAddress(from)
		if err != nil {
			return err
		}

		ucoin := uint64(viper.GetInt(flagAmount))
		tx := transaction.NewTransferTx(from, tos[0], uint64(viper.GetInt(flagGas)), nonce, types.Coin(ucoin), false)

		password := viper.GetString(flagPassword)
		if len(password) < 1 {
			var err error
			password, err = getPassword()
			if err != nil {
				return err
			}
		}

		signedData, err := getSignedDataWithTx(ctx, tx, password, from)
		res, err := ctx.BroadcastSignedData(signedData)
		if err != nil {
			return err
		}
		ctx.PrintOutput(res)
		return nil
	},
}

