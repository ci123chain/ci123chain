package cmd

import (
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/util"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const flagData = "data"

func init()  {
	rootCmd.AddCommand(broadCastCmd)
	broadCastCmd.Flags().String(flagData, "", "signed proposal data")
	util.CheckRequiredFlag(broadCastCmd, flagData)
}

var broadCastCmd = &cobra.Command{
	Use: "broadcast",
	Short: "broadcast transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return err
		}
		data := viper.GetString(flagData)
		dataB, err := hex.DecodeString(data)
		if err != nil {
			return err
		}
		txid, err := ctx.BroadcastSignedData(dataB)

		if err != nil {
			return err
		}
		fmt.Println(txid)
		return nil
	},
}
