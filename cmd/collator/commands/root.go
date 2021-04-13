package commands

import (
	"bufio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	defaultHome = os.ExpandEnv("$HOME/.collator")
	appName     = "clt"
	config      *Config
	homePath    string
	debug       bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}



// NewRootCmd returns the root command for relayer.
func NewRootCmd() *cobra.Command {
	// RootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "This application collator data between configured IBC enabled chains",
		Long: strings.TrimSpace(`The relayer has commands for:
  1. Configuration of the Chains and Paths that the relayer with transfer packets over
  2. Management of keys and light clients on the local machine that will be used to sign and verify txs
  3. Query and transaction functionality for IBC
  4. A responsive relaying application that listens on a path
  5. Commands to assist with development, testnets, and versioning.

NOTE: Most of the commands have aliases that make typing them much quicker (i.e. 'clt tx', 'clt q', etc...)`),
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// reads `homeDir/configs/configs.yaml` into `var configs *Config` before each command
		return initConfig(rootCmd)
	}

	// Register top level flags --home and --debug
	rootCmd.PersistentFlags().StringVar(&homePath, flagHome, defaultHome, "set home directory")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")

	if err := viper.BindPFlag(flagHome, rootCmd.Flags().Lookup(flagHome)); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug")); err != nil {
		panic(err)
	}

	// Register subcommands
	rootCmd.AddCommand(
		configCmd(),
		chainsCmd(),
		//pathsCmd(),
		//flags.LineBreak,
		keysCmd(),
		//lightCmd(),
		//flags.LineBreak,
		//transactionCmd(),
		//queryCmd(),
		//startCmd(),
		//flags.LineBreak,
		//getAPICmd(),
		//flags.LineBreak,
		//devCommand(),
		//testnetsCmd(),
		//getVersionCmd(),
	)

	// This is a bit of a cheat :shushing_face:
	// cdc = codecstd.MakeCodec(simapp.ModuleBasics)
	// appCodec = codecstd.NewAppCodec(cdc)

	return rootCmd
}

// readLineFromBuf reads one line from stdin.
func readStdin() (string, error) {
	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(str), err
}