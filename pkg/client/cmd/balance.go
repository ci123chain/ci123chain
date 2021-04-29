package cmd

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init()  {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().String(util.FlagAddress,"", "address")

	util.CheckRequiredFlag(balanceCmd, util.FlagAddress)
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
		address := viper.GetString(util.FlagAddress)
		addr := sdk.HexToAddress(address)
		v, _, err := ctx.GetBalanceByAddress(addr, false, "")
		if err != nil {
			return err
		}
		fmt.Println(v.String())
		return nil
	},
}
