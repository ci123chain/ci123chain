package cmd

import (
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	acc_type "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"math/big"
)

func AddGenesisAccountCmd(ctx *app.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add-genesis-account [address_or_key_name] [coin]",
		Short: "Add genesis account to genesis.json",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			genFile := config.GenesisFile()
			appState, genDoc, err := app.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			var stakingGenesisState staking.GenesisState
			if _, ok := appState[staking.ModuleName]; !ok{
				return errors.New("unexpected genesisState of staking")
			} else {
				cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState)
			}
			addr := ParseAccAddress(args[0])

			coin, err := ParseCoin(stakingGenesisState.Params.BondDenom, args[1])
			if err != nil {
				return err
			}
			genAcc := acc_type.NewBaseAccountWithAddress(addr)
			//if err := genAcc.Validate(); err != nil {
			//	return err
			//}
			err = genAcc.SetCoins(types.NewCoins(coin))
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
				_ = fmt.Errorf("cannot add account at existing address %v", addr)
			}

			genesisAccounts = append(genesisAccounts, &genAcc)
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

func ParseCoin(denom, amount string) (types.Coin, error) {
	x := new(big.Int)
	x, ok := x.SetString(amount, 10)
	if !ok {
		return types.Coin{}, errors.New("parse coin failed")
	}
	return types.NewCoin(denom, types.NewIntFromBigInt(x)), nil
}

func ParseAccAddress(addr string) types.AccAddress {
	return types.HexToAddress(addr)
}