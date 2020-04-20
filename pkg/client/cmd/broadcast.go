package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/util"
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
