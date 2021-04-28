package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/util"
)

func init()  {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().String(helper.FlagAddress,"", "address")

	util.CheckRequiredFlag(balanceCmd, helper.FlagAddress)
}

var balanceCmd = &cobra.Command{
	Use: "balance",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			panic(err)
		}

		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return  err
		}
		address := viper.GetString(helper.FlagAddress)
		addr := sdk.HexToAddress(address)
		v, _, err, _ := ctx.GetBalanceByAddress(addr, false, "")
		if err != nil {
			return err
		}
		fmt.Println(v.String())
		return nil
	},
}
