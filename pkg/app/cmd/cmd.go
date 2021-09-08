package cmd

import (
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client/lite"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

var LineBreak = &cobra.Command{Run: func(cmd *cobra.Command, args []string) {}}


func AddServerCommands(
	ctx *app.Context,
	cdc *amino.Codec,
	rootCmd *cobra.Command,
	appInit app.AppInit,
	appCreator app.AppCreator,
	appExport 	app.AppExporter,
	)  {
	tendermintCmd := &cobra.Command{
		Use:  "tendermint",
		Short: "Tendermint subcommands",
	}
	tendermintCmd.AddCommand(
		showNodeIDCmd(ctx),
		//showValidatorCmd(ctx),
		//showAddressCmd(ctx),
		//validatorCmd(ctx),
		lite.LiteCmd,
		)

	rootCmd.AddCommand(
		initCmd(ctx, cdc, appInit),
		//createCmd(ctx),
		ExportCmd(appExport, ctx.Config.RootDir),
		startCmd(ctx, appCreator, cdc),
		AddGenesisAccountCmd(ctx, cdc),
		AddGenesisValidatorCmd(ctx, cdc),
		AddGenesisShardCmd(ctx, cdc),
		testnetGenCmd(ctx, cdc, appInit),
		testnetAddCmd(ctx, cdc, appInit),
		bootstrapGenCmd(ctx, cdc, appInit),
		bootstrapAddCmd(ctx, cdc, appInit),
		tendermintCmd,
		LineBreak,
		versionCmd,
		)
}