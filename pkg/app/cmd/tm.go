package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/tendermint/tendermint/p2p"
)

// showNodeIDCmd - ported from Tendermint, dump node ID to stdout
func showNodeIDCmd(ctx *app.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ctx.Config
			nodeKey, err := p2p.LoadNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			fmt.Println(nodeKey.ID())
			return nil
		},
	}
}