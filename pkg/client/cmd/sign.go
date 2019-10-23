package cmd

import (
	"CI123Chain/pkg/client"
	"CI123Chain/pkg/client/helper"
	"CI123Chain/pkg/transaction"
	"CI123Chain/pkg/util"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


func init()  {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String(flagTo, "", "Address sending to")
	signCmd.Flags().Uint(flagAmount, 0, "Amount tbe spent")
	signCmd.Flags().Uint(flagGas, 0, "gas for tx")
	signCmd.Flags().String(helper.FlagAddress, "", "Address to sign with")
	util.CheckRequiredFlag(signCmd, flagAmount)
	util.CheckRequiredFlag(signCmd, flagGas)
}


var signCmd = &cobra.Command{
	Use: "sign",
	Short: "Build, Sign transaction",
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
		tx := &transaction.TransferTx{
			Common: transaction.CommonTx{
				Code: transaction.TRANSFER,
				From: from,
				Gas:  uint64(viper.GetInt(flagGas)),
				Nonce:nonce,
			},
			To: tos[0],
			Amount: uint64(viper.GetInt(flagAmount)),
		}
		signedtx, err := ctx.SignTx(tx, from)

		txByte := signedtx.Bytes()
		if err != nil {
			return err
		}
		fmt.Println(hex.EncodeToString(txByte))
		return nil
	},
}

