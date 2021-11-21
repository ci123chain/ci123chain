package cmd

import (
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/app/types"
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
			var err error;
			if validatorKeystr := viper.GetString(app.FlagValidatorKey); len(validatorKeystr) != 0 {
				validatorKey, err = app.CreatePVWithKey(types.GetCodec(), validatorKeystr)
				if err != nil {
					panic(err)
				}
			}
			pv := validator.GenFilePV(
				ctxConfig.PrivValidatorKeyFile(),
				ctxConfig.PrivValidatorStateFile(),
				validatorKey,
			)
			_, _ = node.GenNodeKeyByPrivKey(ctxConfig.NodeKeyFile(), pv.Key.PrivKey)
		},
	}
	cmd.Flags().String(app.FlagValidatorKey, "", "the validator key")
	return cmd
}

