package cmd

import (
	"CI123Chain/pkg/client/helper"
	"CI123Chain/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"CI123Chain/pkg/client"
)

func init()  {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().String(helper.FlagAddress, "", "address")
	util.CheckRequiredFlag(balanceCmd, helper.FlagAddress)
}

var balanceCmd = &cobra.Command{
	Use: "balance",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx, err := client.NewClientContextFromViper()
		if err != nil {
			return  err
		}
		addrs, err := ctx.GetInputAddresses()
		if err != nil {
			return err
		}
		v, err := ctx.GetBalanceByAddress(addrs[0])
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	},
}
