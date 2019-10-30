package cmd

import (
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/app"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/client/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var homeDir = os.ExpandEnv("$HOME/.ci123_client")
var cdc = app.GetCodec()

var rootCmd = &cobra.Command{
	Use: 	"cli", 
	Short:  "Blockchain Client",
}

func Execute()  {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init()  {
	rootCmd.PersistentFlags().StringP(helper.FlagHomeDir, "", homeDir, "directory for keystore")
	rootCmd.PersistentFlags().BoolP(helper.FlagVerbose, "v", false, "enable verbose output")
	rootCmd.PersistentFlags().String(helper.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	rootCmd.PersistentFlags().StringP(helper.FlagPassword, "p", "", "password for signing tx")
	rootCmd.PersistentFlags().Int64(helper.FlagHeight, 0, "Use a special height to query state at (this can error if node is pruning state)")
	viper.BindPFlags(rootCmd.Flags())
}