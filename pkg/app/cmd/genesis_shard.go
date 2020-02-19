package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/order"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	otype "github.com/tanhuiya/ci123chain/pkg/order/types"
	"github.com/tendermint/tendermint/libs/cli"
	"strconv"
	"strings"
)

func AddGenesisShardCmd(ctx *app.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add-genesis-shard [name1:height1;name2:height2]",
		Short: "Add genesis shard to genesis.json",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			var names []string
			var heights []int64
			shards := strings.Split(args[0],";")
			for i := 0; i < len(shards); i++ {
				shard := strings.Split(shards[i],":")
				names = append(names,shard[0])
				height, _ := strconv.ParseInt(shard[1],10,64)
				heights = append(heights, height)
			}
			genFile := config.GenesisFile()
			appState, genDoc, err := app.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}
			var gs otype.GenesisState
			if _, ok := appState[order.ModuleName]; !ok {
				gs = otype.GenesisState{}
			} else {
				cdc.MustUnmarshalJSON(appState[order.ModuleName], &gs)
			}
			if gs.Params.OrderBook.Lists[0].Name == ""{
				gs.Params.OrderBook.Lists = []keeper.Lists{}
			}
			for i := 0; i < len(shards); i++ {
				var list keeper.Lists
				exist := false
				list.Name = names[i]
				list.Height = heights[i]
				for _,v := range gs.Params.OrderBook.Lists {
					if v.Name == list.Name{
						exist = true
						break
					}
				}
				if exist {
					continue
				}
				gs.Params.OrderBook.Lists = append(gs.Params.OrderBook.Lists, list)
			}
			gBz := cdc.MustMarshalJSON(gs)
			appState[order.ModuleName] = gBz
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}
			genDoc.AppState = appStateJSON
			return app.ExportGenesisFile(genDoc, genFile)
		},
	}
	return cmd
}