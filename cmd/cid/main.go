package main

import (
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/app/cmd"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"os"
)

const (

	DefaultConfDir = "$HOME/.ci123"
	flagLogLevel = "log_level"
	//logDEBUG     = "main:debug,state:debug,ibc:debug,*:error"
	logINFO      = "x/upgrade:info,x/registry:info,state:debug,x/staking:info,x/preStaking:info,x/ibc/client:info,x/ibc/connection:info,x/ibc/channel:info,x/gravity:info,x/supply:info,*:error"
	logDEBUG      = "*:debug"
	//logINFO      = "*:info"
	logERROR     = "*:error"
	logNONE      = "*:none"
)


func main()  {
	cobra.EnableCommandSorting = false
	ctx := app.NewDefaultContext()
	app.GasPriceConfig = ctx.Config.GasPrice
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
	rootCmd.Flags().String(app.HomeFlag, DefaultConfDir, "directory for configs and data")
	rootCmd.Flags().String(flagLogLevel, "info", "Run abci app with different log level")
	rootCmd.PersistentFlags().String("log_level", ctx.Config.LogLevel, "log level")

	rootDir := os.ExpandEnv(viper.GetString(app.HomeFlag))
	ctx.Config.SetRoot(rootDir)
	cmd.AddServerCommands(
		ctx,
		types2.GetCodec(),
		rootCmd,
		app.NewAppInit(),
		app.ConstructAppCreator(app.NewApp, app.AppName),
		app.ConstructAppExporter(app.AppName),
		)
	viper.SetEnvPrefix("CI")
	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.AutomaticEnv()
	//rootDir := os.ExpandEnv(DefaultConfDir)
	//if len(viper.GetString(HomeFlag)) > 0 {
	//	rootDir = os.ExpandEnv(viper.GetString(HomeFlag))
	//}
	exector := cli.PrepareBaseCmd(rootCmd, "CI", rootDir)
	exector.Execute()
}


//func exportAppState(lg log.Logger, ldb db.DB, cdb db.DB, traceStore io.Writer) (json.RawMessage, []types.GenesisValidator, error) {
//	logger.SetLogger(lg)
//	return app.NewChain(lg, ldb, cdb, traceStore).ExportAppStateJSON()
//}