package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	distr "github.com/ci123chain/ci123chain/pkg/distribution/types"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"strconv"
	"time"
)


func AddGenesisValidatorCmd(ctx *app.Context, cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:  "add-genesis-validator [address] [amount] [pub_key] [commission_rate] [commission_max_rate] [commission_max_change_rate]",
		Short: "Add genesis validator to genesis.json",
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr := ParseAccAddress(args[0])

			coin, err := parseAmount(args[1])
			if err != nil {
				return err
			}

			_, err = parsePubKey(args[2], cdc)
			if err != nil {
				return err
			}
			timeNow := time.Now()
			commission, err := parseCommission(args[3], args[4], args[5], timeNow)
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := app.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}
			var stakingGenesisState staking.GenesisState
			var genesisValidator staking.Validator
			if _, ok := appState[staking.ModuleName]; !ok{
				stakingGenesisState = staking.GenesisState{}
			} else {
				cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState)
			}
			genesisValidator = staking.Validator{
				OperatorAddress:   addr,
				ConsensusKey:      args[2],

				Jailed:            false,
				Status:            1,
				Tokens:            types.NewInt(coin),
				DelegatorShares:   types.NewDec(coin),
				Description:       staking.Description{},
				UnbondingHeight:   0,
				UnbondingTime:     time.Time{},
				Commission:        commission,
				MinSelfDelegation: types.NewInt(coin),
			}

			delegation := staking.NewDelegation(addr, addr, types.NewDec(coin))

			stakingGenesisState.Validators = append(stakingGenesisState.Validators, genesisValidator)
			stakingGenesisState.Delegations = append(stakingGenesisState.Delegations, delegation)

			genesisStateBz := cdc.MustMarshalJSON(stakingGenesisState)
			appState[staking.ModuleName] = genesisStateBz

			//distribution
			var distributionGenesisState distr.GenesisState
			if _, ok := appState[distr.ModuleName]; !ok {
				distributionGenesisState = distr.GenesisState{}
			}else {
				cdc.MustUnmarshalJSON(appState[distr.ModuleName], &distributionGenesisState)
			}
			outstanddingReward := distr.ValidatorOutstandingRewardsRecord{
				ValidatorAddress:   addr,
				OutstandingRewards: types.NewEmptyDecCoin(),
			}
			currentReward := distr.ValidatorCurrentRewardsRecord{
				ValidatorAddress: addr,
				Rewards:          distr.ValidatorCurrentRewards{
					Rewards: types.NewEmptyDecCoin(),
					Period:  0,
				},
			}
			distributionGenesisState.ValidatorCurrentRewards = append(distributionGenesisState.ValidatorCurrentRewards, currentReward)
			distributionGenesisState.OutstandingRewards = append(distributionGenesisState.OutstandingRewards, outstanddingReward)
			distrGenesisStateBz := cdc.MustMarshalJSON(distributionGenesisState)
			appState[distr.ModuleName] = distrGenesisStateBz

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

func parseAmount(amt string) (int64, error) {
	amount, err := strconv.ParseInt(amt, 10, 64)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func parsePubKey(pub string, cdc *codec.Codec) (crypto.PubKey, error) {
	pubByte, err := hex.DecodeString(pub)
	if err != nil {
		return nil, err
	}
	var public crypto.PubKey
	err = cdc.UnmarshalJSON(pubByte, &public)
	if err != nil {
		return nil, err
	}
	return public, nil
}

func parseCommission(rateStr, maxRateStr, maxChangeRateStr string, timeNow time.Time) (staking.Commission, error) {
	rate, err := strconv.ParseInt(rateStr, 10, 64)
	if err != nil {
		return staking.Commission{}, err
	}
	maxRae, err := strconv.ParseInt(maxRateStr, 10, 64)
	if err != nil {
		return staking.Commission{}, err
	}
	maxChangeRate, err := strconv.ParseInt(maxChangeRateStr, 10, 64)
	if err != nil {
		return staking.Commission{}, err
	}
	err = checkParam(rate, maxRae, maxChangeRate)
	if err != nil {
		return staking.Commission{}, err
	}
	return staking.Commission{
		CommissionRates: staking.CommissionRates{
			Rate:          types.NewDecWithPrec(rate, 2),
			MaxRate:       types.NewDecWithPrec(maxRae, 2),
			MaxChangeRate: types.NewDecWithPrec(maxChangeRate, 2),
		},
		UpdateTime:      timeNow,
	}, nil
}

func checkParam(keys... int64) error {
	for _, k := range keys {
		if k > 100 || k <=0 {
			return errors.New(fmt.Sprintf("invalid params %d", k))
		}
	}
	return nil
}