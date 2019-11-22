package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/util"
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
		ctx, err := client.NewClientContextFromViper()
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
