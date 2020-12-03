package cmd

import (
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	distr "github.com/ci123chain/ci123chain/pkg/distribution/types"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

			coin, err := util.CheckInt64(args[1])
			if err != nil {
				return err
			}

			_, err = util.ParsePubKey(args[2])
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
				UnbondingHeight:   int64(-1),
				UnbondingTime:     time.Time{},
				BondedHeight:     int64(0),
				Commission:        commission,
				MinSelfDelegation: types.NewInt(coin),
			}
			genesisValidator.ConsensusAddress = genesisValidator.GetConsPubKey().Address().String()

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

func parseCommission(rateStr, maxRateStr, maxChangeRateStr string, timeNow time.Time) (staking.Commission, error) {
	rate, err := strconv.ParseInt(rateStr, 10, 64)
	if err != nil {
		return staking.Commission{}, err
	}
	maxRate, err := strconv.ParseInt(maxRateStr, 10, 64)
	if err != nil {
		return staking.Commission{}, err
	}
	maxChangeRate, err := strconv.ParseInt(maxChangeRateStr, 10, 64)
	if err != nil {
		return staking.Commission{}, err
	}
	err = checkParam(rate, maxRate, maxChangeRate)
	if err != nil {
		return staking.Commission{}, err
	}
	if rate > maxRate {
		return staking.Commission{}, errors.New("rate can't grater than max_rate")
	}
	if maxChangeRate > maxRate {
		return staking.Commission{}, errors.New("max_change_rate can't grater than max_rate")
	}
	return staking.Commission{
		CommissionRates: staking.CommissionRates{
			Rate:          types.NewDecWithPrec(rate, 2),
			MaxRate:       types.NewDecWithPrec(maxRate, 2),
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