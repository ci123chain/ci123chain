package cmd

import (
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/node"
	"github.com/ci123chain/ci123chain/pkg/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

func genValidatorCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "gen-validator",
		Short: "Generate new validator keypair",
		Run: func(cmd *cobra.Command, args []string) {
			ctxConfig := ctx.Config
			ctxConfig.SetRoot(viper.GetString(tmcli.HomeFlag))

			validatorKey := ed25519.GenPrivKey()
			pv := validator.GenFilePV(
				ctxConfig.PrivValidatorKeyFile(),
				ctxConfig.PrivValidatorStateFile(),
				validatorKey,
			)
			_, _ = node.GenNodeKeyByPrivKey(ctxConfig.NodeKeyFile(), pv.Key.PrivKey)
		},
	}
	return cmd
}

