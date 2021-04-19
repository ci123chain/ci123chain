package commands

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	"github.com/spf13/cobra"
	"strings"
)

// transactionCmd returns a parent transaction command handler, where all child
// commands can submit transactions on IBC-connected networks.
func transactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transact",
		Aliases: []string{"tx"},
		Short:   "IBC transaction commands",
		Long: strings.TrimSpace(`Commands to create IBC transactions on pre-configured chains.
Most of these commands take a [path] argument. Make sure:
  1. Chains are properly configured to relay over by using the 'rly chains list' command
  2. Path is properly configured to relay over by using the 'rly paths list' command`,
		),
	}

	cmd.AddCommand(
		linkCmd(),
		//linkThenStartCmd(),
		//relayMsgsCmd(),
		//relayAcksCmd(),
		//xfersend(),
		//flags.LineBreak,
		//createClientsCmd(),
		//updateClientsCmd(),
		//upgradeClientsCmd(),
		//upgradeChainCmd(),
		//createConnectionCmd(),
		//closeChannelCmd(),
		//flags.LineBreak,
		//rawTransactionCmd(),
		//flags.LineBreak,
		//sendCmd(),
	)

	return cmd
}

func linkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "link [path-name]",
		Aliases: []string{"connect"},
		Short:   "create clients, connection, and channel between two configured chains with a configured path",
		Long: strings.TrimSpace(`Create an IBC client between two IBC-enabled networks, in addition
to creating a connection and a channel between the two networks on a configured path.`,
		),
		Args: cobra.ExactArgs(1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s transact link demo-path
$ %s tx link demo-path
$ %s tx connect demo-path`,
			appName, appName, appName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			to, err := getTimeout(cmd)
			if err != nil {
				return err
			}

			retries, err := cmd.Flags().GetUint64(flagMaxRetries)
			if err != nil {
				return err
			}

			// ensure that keys exist
			if _, err = c[src].GetAddress(); err != nil {
				return err
			}
			if _, err = c[dst].GetAddress(); err != nil {
				return err
			}

			// create clients if they aren't already created
			modified, err := c[src].CreateClients(c[dst])
			if modified {
				if err := overWriteConfig(config); err != nil {
					return err
				}
			}

			if err != nil {
				return err
			}

			// create connection if it isn't already created
			modified, err = c[src].CreateOpenConnections(c[dst], retries, to)
			if modified {
				if err := overWriteConfig(config); err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}

			// create channel if it isn't already created
			modified, err = c[src].CreateOpenChannels(c[dst], 3, to)
			if modified {
				if err := overWriteConfig(config); err != nil {
					return err
				}
			}
			return err
		},
	}

	return retryFlag(timeoutFlag(cmd))
}
// ensureKeysExist returns an error if a configured key for a given chain does
// not exist.
func ensureKeysExist(chains map[string]*collactor.Chain) error {
	for _, v := range chains {
		if _, err := v.GetAddress(); err != nil {
			return err
		}
	}

	return nil
}
