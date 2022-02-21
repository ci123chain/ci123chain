package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"gopkg.in/yaml.v2"
)

const flagShort = "short"

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the app version",
		RunE:   printVersion,
	}
	cmd.Flags().Bool(flagShort, false, "Print long version information")
	return cmd
}

// return version of CLI/node and commit hash
func GetVersion() version.Info {
	verInfo := version.NewInfo()
	return verInfo
}

// CMD
func printVersion(cmd *cobra.Command, args []string) error {
	verInfo := GetVersion()
	if viper.GetBool(flagShort) {
		fmt.Println(verInfo.Version)
		return nil
	}

	var bz []byte
	var err error

	switch viper.GetString(cli.OutputFlag) {
	case "json":
		bz, err = json.Marshal(verInfo)
	default:
		bz, err = yaml.Marshal(&verInfo)
	}

	if err != nil {
		return err
	}

	_, err = fmt.Println(string(bz))
	return err
}
