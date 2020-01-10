package main

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/app/cmd"
	"github.com/tanhuiya/ci123chain/pkg/logger"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tm-db"
	"io"
	"os"
)

const (
	appName = "ci123"
	DefaultConfDir = "$HOME/.ci123"
	flagLogLevel       = "log_level"
	logDEBUG     = "main:debug,state:debug,ibc:debug,*:error"
	logINFO      = "main:info,state:info,ibc:info,*:error"
	logERROR     = "*:error"
	logNONE      = "*:none"
)

func main()  {
	cobra.EnableCommandSorting = false
	ctx := app.NewDefaultContext()
	rootCmd := &cobra.Command{
		Use: 	"ci123",
		Short:  "ci123 node",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch viper.GetString(flagLogLevel) {
			case "debug":
				return app.SetupContext(ctx, logDEBUG)
			case "info":
				return app.SetupContext(ctx, logINFO)
			case "error":
				return app.SetupContext(ctx, logERROR)
			case "none":
				return app.SetupContext(ctx, logNONE)
			}
			return app.SetupContext(ctx, logINFO)
		},
	}
	rootCmd.Flags().String(flagLogLevel, "info", "Run abci app with different log level")
	rootCmd.PersistentFlags().String("log_level", ctx.Config.LogLevel, "log level")

	cmd.AddServerCommands(
		ctx,
		app.MakeCodec(),
		rootCmd,
		app.NewAppInit(),
		app.ConstructAppCreator(newApp, appName),
		app.ConstructAppExporter(exportAppState, appName),
		)

	viper.BindPFlags(rootCmd.Flags())
	rootDir := os.ExpandEnv(DefaultConfDir)
	if len(viper.GetString(cli.HomeFlag)) > 0 {
		rootDir = os.ExpandEnv(viper.GetString(cli.HomeFlag))
	}
	exector := cli.PrepareBaseCmd(rootCmd, "CORE", rootDir)
	exector.Execute()
}

func newApp(lg log.Logger, db db.DB, traceStore io.Writer) abci.Application{
	logger.SetLogger(lg)
	//将ibc的logger设置为ibc module.
	return app.NewChain(lg.With("module", "ibc"), db, traceStore)
}

func exportAppState(lg log.Logger, db db.DB, traceStore io.Writer) (json.RawMessage, []types.GenesisValidator, error) {
	logger.SetLogger(lg)
	return app.NewChain(lg, db, traceStore).ExportAppStateJSON()
}