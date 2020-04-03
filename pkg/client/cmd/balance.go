package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/util"
)

func init()  {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().String(helper.FlagAddress,"", "address")
	err := viper.BindPFlags(balanceCmd.Flags())
	if err != nil {
		panic(err)
	}
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
		address := viper.GetString(helper.FlagAddress)
		addr := sdk.HexToAddress(address)
		v, err := ctx.GetBalanceByAddress(addr)
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	},
}
