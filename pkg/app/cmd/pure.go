package cmd

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

const flagPureHeight = "pure_height"

func pureStateCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prueState",
		Short: "prue state latest state",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- prue start ---------")
			pureHeight := viper.GetInt(flagPureHeight)
			if pureHeight < 1 {
				fmt.Println("pure_height not provided")
			}
			pureState(ctx, pureHeight)
			log.Println("--------- replay success ---------")
		},
	}
	cmd.Flags().IntP(flagPureHeight, "", 0, "height for stop replaying")
	return cmd
}

func pureState(ctx *app.Context, pureHeight int)  {

}