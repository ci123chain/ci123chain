package cmd

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init()  {
	rootCmd.AddCommand(nonceCmd)
	nonceCmd.Flags().String(helper.FlagAddress, "", "address")
	util.CheckRequiredFlag(nonceCmd, helper.FlagAddress)
}

var nonceCmd = &cobra.Command{
	Use: "nonce",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return  err
		}
		addr := ctx.GetFromAddresses()
		v, _, err := ctx.GetNonceByAddress(addr, false)
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	},
}
