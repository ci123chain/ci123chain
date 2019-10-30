package cmd

import (
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/client/helper"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/client"
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
