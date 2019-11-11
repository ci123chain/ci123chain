package main

import (
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/app/cmd"
	"github.com/tanhuiya/ci123chain/pkg/logger"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tm-db"
	"github.com/tendermint/tendermint/libs/log"
	abci "github.com/tendermint/tendermint/abci/types"
	"io"
	"os"
)

const (
	appName = "ci123"
	confDir = "$HOME/.ci123"
)

func main()  {
	cobra.EnableCommandSorting = false
	ctx := app.NewDefaultContext()
	rootCmd := &cobra.Command{
		Use: 	"ci123",
		Short:  "ci123 node",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.SetupContext(ctx)
		},
	}

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
	rootDir := os.ExpandEnv(confDir)
	exector := cli.PrepareBaseCmd(rootCmd, "PC", rootDir)
	exector.Execute()
}

func newApp(lg log.Logger, db db.DB, traceStore io.Writer) abci.Application{
	logger.SetLogger(lg)
	return app.NewChain(lg, db, traceStore)
}

func exportAppState(lg log.Logger, db db.DB, traceStore io.Writer) (json.RawMessage, []types.GenesisValidator, error) {
	logger.SetLogger(lg)
	return app.NewChain(lg, db, traceStore).ExportAppStateJSON()
}