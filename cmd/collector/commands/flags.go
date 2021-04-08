package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

const (
	flagHome = "home"
	flagForce               = "force"
	flagTimeout             = "timeout"
	flagFile                = "file"
	flagURL                 = "url"
	flagMaxRetries          = "max-retries"

)


func getAddInputs(cmd *cobra.Command) (file string, url string, err error) {
	file, err = cmd.Flags().GetString(flagFile)
	if err != nil {
		return
	}

	url, err = cmd.Flags().GetString(flagURL)
	if err != nil {
		return
	}

	if file != "" && url != "" {
		return "", "", errMultipleAddFlags
	}

	return
}


func getTimeout(cmd *cobra.Command) (time.Duration, error) {
	to, err := cmd.Flags().GetString(flagTimeout)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(to)
}

func chainsAddFlags(cmd *cobra.Command) *cobra.Command {
	fileFlag(cmd)
	urlFlag(cmd)
	return cmd
}

func fileFlag(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP(flagFile, "f", "", "fetch json data from specified file")
	if err := viper.BindPFlag(flagFile, cmd.Flags().Lookup(flagFile)); err != nil {
		panic(err)
	}
	return cmd
}


func urlFlag(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP(flagURL, "u", "", "url to fetch data from")
	if err := viper.BindPFlag(flagURL, cmd.Flags().Lookup(flagURL)); err != nil {
		panic(err)
	}
	return cmd
}



func forceFlag(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP(flagForce, "f", false,
		"option to force non-standard behavior such as initialization of light client from"+
			"configured chain or generation of new path")
	if err := viper.BindPFlag(flagForce, cmd.Flags().Lookup(flagForce)); err != nil {
		panic(err)
	}
	return cmd
}
