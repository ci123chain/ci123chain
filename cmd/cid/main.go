package main

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/app/cmd"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	flagLogLevel = "log_level"
	HomeFlag     = "home"
	//logDEBUG     = "main:debug,state:debug,ibc:debug,*:error"
	logINFO      = "state:info,x/ibc/client:info,x/ibc/connection:info,x/ibc/channel:info,*:error"
	logDEBUG      = "*:debug"
	//logINFO      = "*:info"
	logERROR     = "*:error"
	logNONE      = "*:none"
)

func main()  {
	cobra.EnableCommandSorting = false
	ctx := app.NewDefaultContext()
	rootCmd := &cobra.Command{
		Use: 	 "ci123",
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
	rootCmd.Flags().String(HomeFlag, "", "directory for configs and data")
	rootCmd.Flags().String(flagLogLevel, "info", "Run abci app with different log level")
	rootCmd.PersistentFlags().String("log_level", ctx.Config.LogLevel, "log level")
	rootCmd.Flags().String(app.FlagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")

	cmd.AddServerCommands(
		ctx,
		types2.GetCodec(),
		rootCmd,
		app.NewAppInit(),
		app.ConstructAppCreator(newApp, appName),
		app.ConstructAppExporter(exportAppState, appName),
		)
	viper.SetEnvPrefix("CI")
	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.AutomaticEnv()
	rootDir := os.ExpandEnv(DefaultConfDir)
	if len(viper.GetString(HomeFlag)) > 0 {
		rootDir = os.ExpandEnv(viper.GetString(HomeFlag))
	}
	exector := cli.PrepareBaseCmd(rootCmd, "CI", rootDir)
	exector.Execute()
}

func newApp(lg log.Logger, ldb db.DB, cdb db.DB,traceStore io.Writer) abci.Application{
	logger.SetLogger(lg)
	return app.NewChain(lg, ldb, cdb, traceStore)
}

func exportAppState(lg log.Logger, ldb db.DB, cdb db.DB, traceStore io.Writer) (json.RawMessage, []types.GenesisValidator, error) {
	logger.SetLogger(lg)
	return app.NewChain(lg, ldb, cdb, traceStore).ExportAppStateJSON()
}