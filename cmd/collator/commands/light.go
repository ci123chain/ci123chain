package commands

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/collactor/helper"
	"github.com/spf13/cobra"
	"strings"
)

// chainCmd represents the keys command
func lightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "light",
		Aliases: []string{"l"},
		Short:   "manage light clients held by the relayer for each chain",
	}

	//cmd.AddCommand(lightHeaderCmd())
	cmd.AddCommand(initLightCmd())
	//cmd.AddCommand(updateLightCmd())
	//cmd.AddCommand(deleteLightCmd())

	return cmd
}


func initLightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [chain-id]",
		Aliases: []string{"i"},
		Short:   "Initiate the light client",
		Long: `Initiate the light client by:
	1. passing it a root of trust as a --hash/-x and --height
	2. Use --force/-f to initialize from the configured node`,
		Args: cobra.ExactArgs(1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s light init ibc-0 --force
$ %s light init ibc-1 --height 1406 --hash <hash>
$ %s l i ibc-2 --force`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}

			force, err := cmd.Flags().GetBool(flagForce)
			if err != nil {
				return err
			}
			//height, err := cmd.Flags().GetInt64(flags.FlagHeight)
			//if err != nil {
			//	return err
			//}
			//hash, err := cmd.Flags().GetBytesHex(flagHash)
			//if err != nil {
			//	return err
			//}

			out, err := helper.InitLight(chain, force)
			if err != nil {
				return err
			}

			if out != "" {
				fmt.Println(out)
			}

			return nil
		},
	}
	return forceFlag(cmd)
	//return forceFlag(lightFlags(cmd))

}