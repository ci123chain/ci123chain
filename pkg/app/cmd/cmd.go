package cmd

import (
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"

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
		tcmd.LiteCmd,
		)

	rootCmd.AddCommand(
		initCmd(ctx, cdc, appInit),
		//createCmd(ctx),
		startCmd(ctx, appCreator),
		AddGenesisAccountCmd(ctx, cdc),
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