package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/util"
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
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return  err
		}
		addrs, err := ctx.GetInputAddresses()
		if err != nil {
			return err
		}
		v, err := ctx.GetNonceByAddress(addrs[0])
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	},
}
