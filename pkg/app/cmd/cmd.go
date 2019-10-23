package cmd

import (
	"CI123Chain/pkg/app"
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
		//showNodeIDCmd(ctx),
		//showValidatorCmd(ctx),
		//showAddressCmd(ctx),
		//validatorCmd(ctx),
		)

	rootCmd.AddCommand(
		initCmd(ctx, cdc, appInit),
		//createCmd(ctx),
		startCmd(ctx, appCreator),
		tendermintCmd,
		LineBreak,
		versionCmd,
		)
}