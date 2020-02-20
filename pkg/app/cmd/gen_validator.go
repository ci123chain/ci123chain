package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tanhuiya/ci123chain/pkg/app"
)

func genValidatorCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "gen-validator",
		Short: "Generate new validator keypair",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmd
}

