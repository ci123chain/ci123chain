package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	acc_type "github.com/tanhuiya/ci123chain/pkg/account/types"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tendermint/tendermint/libs/cli"
	"strconv"
)

func AddGenesisAccountCmd(ctx *app.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add-genesis-account [address_or_key_name] [coin]",
		Short: "Add genesis account to genesis.json",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr := ParseAccAddress(args[0])

			coin, err := ParseCoin(args[1])
			if err != nil {
				return err
			}
			genAcc := account.NewGenesisAccountRaw(addr, coin)
			if err := genAcc.Validate(); err != nil {
				return err
			}
			genFile := config.GenesisFile()
			appState, genDoc, err := app.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			var genesisAccounts acc_type.GenesisAccounts
			if _, ok := appState[account.ModuleName]; !ok {
				genesisAccounts = acc_type.GenesisAccounts{}
			} else {
				cdc.MustUnmarshalJSON(appState[account.ModuleName], &genesisAccounts)
			}
			if genesisAccounts.Contains(addr) {
				fmt.Errorf("cannot add account at existing address %v", addr)
			}

			genesisAccounts = append(genesisAccounts, genAcc)
			genesisStateBz := cdc.MustMarshalJSON(account.GenesisState(genesisAccounts))
			appState[account.ModuleName] = genesisStateBz

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

func ParseCoin(coin string) (types.Coin, error) {
	coin64, err := strconv.ParseUint(coin, 10, 64)
	return types.NewUInt64Coin(coin64), err
}

func ParseAccAddress(addr string) types.AccAddress {
	return types.HexToAddress(addr)
}