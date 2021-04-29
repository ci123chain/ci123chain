package main

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/app/cmd"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ci123chain/ci123chain/pkg/util"
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

//const (
//	appName = "ci123"
//	DefaultConfDir = "$HOME/.ci123"
//	flagLogLevel = "log_level"
//	HomeFlag     = "home"
//	//logDEBUG     = "main:debug,state:debug,ibc:debug,*:error"
//	logINFO      = "state:info,x/ibc/client:info,x/ibc/connection:info,x/ibc/channel:info,*:error"
//	logDEBUG      = "*:debug"
//	//logINFO      = "*:info"
//	logERROR     = "*:error"
//	logNONE      = "*:none"
//)

func main()  {
	cobra.EnableCommandSorting = false
	ctx := app.NewDefaultContext()
	rootCmd := &cobra.Command{
		Use: 	 "ci123",
		Short:  "ci123 node",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch viper.GetString(util.FlagLogLevel) {
			case "debug":
				return app.SetupContext(ctx, util.LogDEBUG)
			case "info":
				return app.SetupContext(ctx, util.LogINFO)
			case "error":
				return app.SetupContext(ctx, util.LogERROR)
			case "none":
				return app.SetupContext(ctx, util.LogNONE)
			}
			return app.SetupContext(ctx, util.LogINFO)
		},
	}
	rootCmd.Flags().String(util.HomeFlag, "", "directory for configs and data")
	//rootCmd.Flags().String(util.FlagLogLevel, "info", "Run abci app with different log level")
	rootCmd.PersistentFlags().String(util.FlagLogLevel, ctx.Config.LogLevel, "log level")
	rootCmd.Flags().String(util.FlagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")

	cmd.AddServerCommands(
		ctx,
		types2.GetCodec(),
		rootCmd,
		app.NewAppInit(),
		app.ConstructAppCreator(newApp, util.AppName),
		app.ConstructAppExporter(exportAppState, util.AppName),
		)
	viper.SetEnvPrefix("CI")
	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.AutomaticEnv()
	rootDir := os.ExpandEnv(util.DefaultConfDir)
	if len(viper.GetString(util.HomeFlag)) > 0 {
		rootDir = os.ExpandEnv(viper.GetString(util.HomeFlag))
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