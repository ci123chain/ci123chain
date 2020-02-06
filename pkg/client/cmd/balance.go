package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/util"
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
		ctx, err := client.NewClientContextFromViper(cdc)
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
