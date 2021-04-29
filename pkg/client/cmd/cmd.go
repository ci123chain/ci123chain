package cmd

import (
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var homeDir = os.ExpandEnv("$HOME/.cli")
var cdc = types.GetCodec()

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
	rootCmd.PersistentFlags().StringP(util.FlagHomeDir, "", homeDir, "directory for keystore")
	rootCmd.PersistentFlags().BoolP(util.FlagVerbose, "v", false, "enable verbose output")
	rootCmd.PersistentFlags().String(util.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	rootCmd.PersistentFlags().StringP(util.FlagPassword, "p", "", "password for signing tx")
	rootCmd.PersistentFlags().Int64(util.FlagHeight, 0, "Use a special height to query state at (this can error if node is pruning state)")
	viper.SetEnvPrefix("CI")
	_ = viper.BindPFlags(rootCmd.Flags())
	_ = viper.BindPFlags(rootCmd.PersistentFlags())
	viper.AutomaticEnv()
}