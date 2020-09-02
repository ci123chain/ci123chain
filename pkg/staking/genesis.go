package staking

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	account "github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/supply/exported"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis sets the pool and parameters for the provided keeper.  For each
// validator in data, it sets that validator in the keeper along with manually
// setting the indexes. In addition, it also sets any delegations found in
// data. Finally, it updates the bonded validators.
// Returns final validator set after applying all declaration and delegations
func InitGenesis(
	ctx sdk.Context, keeper keeper.StakingKeeper, ak account.AccountKeeper,
	supplyKeeper supply.Keeper, data types.GenesisState,
) (res []abci.ValidatorUpdate) {

	bondedTokens := sdk.ZeroInt()
	notBondedTokens := sdk.ZeroInt()

	// We need to pretend to be "n blocks before genesis", where "n" is the
	// validator update delay, so that e.g. slashing periods are correctly
	// initialized for the validator set e.g. with a one-block offset - the
	// first TM block is at height 1, so state updates applied from
	// genesis.json are in block 0.
	ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)

	keeper.SetParams(ctx, data.Params)
	keeper.SetLastTotalPower(ctx, data.LastTotalPower)

	for _, validator := range data.Validators {
		//opAddr := sdk.ToAccAddress(validator.GetConsPubKey().Address())
		//validator.OperatorAddress = opAddr
		err := keeper.SetValidator(ctx, validator)
		if err != nil {
			panic(err)
		}

		// Manually set indices for the first time
		keeper.SetValidatorByConsAddr(ctx, validator)
		keeper.SetValidatorByPowerIndex(ctx, validator)

		// Call the creation hook if not exported
		if !data.Exported {
			keeper.AfterValidatorCreated(ctx, validator.OperatorAddress)
		}

		// update timeslice if necessary
		if validator.IsUnbonding() {
			keeper.InsertValidatorQueue(ctx, validator)
		}

		switch validator.GetStatus() {
		case sdk.Bonded:
			bondedTokens = bondedTokens.Add(validator.GetTokens())
		case sdk.Unbonding, sdk.Unbonded:
			notBondedTokens = notBondedTokens.Add(validator.GetTokens())
		default:
			panic("invalid validator status")
		}
	}

	for _, delegation := range data.Delegations {
		// Call the before-creation hook if not exported
		if !data.Exported {
			keeper.BeforeDelegationCreated(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
		}
		keeper.SetDelegation(ctx, delegation)

		// Call the after-modification hook if not exported
		if !data.Exported {
			keeper.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
		}
	}

	for _, ubd := range data.UnbondingDelegations {
		keeper.SetUnbondingDelegation(ctx, ubd)
		for _, entry := range ubd.Entries {
			keeper.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
			notBondedTokens = notBondedTokens.Add(entry.Balance)
		}
	}

	for _, red := range data.Redelegations {
		keeper.SetRedelegation(ctx, red)
		for _, entry := range red.Entries {
			keeper.InsertRedelegationQueue(ctx, red, entry.CompletionTime)
		}
	}

	//bondedCoins := sdk.NewCoins(sdk.NewCoin(bondedTokens))
	//notBondedCoins := sdk.NewCoins(sdk.NewCoin(notBondedTokens))

	// check if the unbonded and bonded pools accounts exists
	bondedPool := keeper.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	// add coins if not provided on genesis
	if ak.GetAllBalances(ctx, bondedPool.GetAddress()).IsZero() {
		ModuleAcc := ak.GetAccount(ctx, bondedPool.GetAddress()).(exported.ModuleAccountI)
		err := ModuleAcc.SetCoin(sdk.NewCoin(bondedTokens))
		if err != nil {
			panic(err)
		}

		supplyKeeper.SetModuleAccount(ctx, ModuleAcc)
	}

	notBondedPool := keeper.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	if ak.GetAllBalances(ctx, notBondedPool.GetAddress()).IsZero() {
		ModuleAcc := ak.GetAccount(ctx, notBondedPool.GetAddress()).(exported.ModuleAccountI)
		err := ModuleAcc.SetCoin(sdk.NewCoin(notBondedTokens))
		if err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, ModuleAcc)
	}
	//var coins = supplyKeeper.GetModuleAccount(ctx, notBondedPool.GetName()).GetCoin()
	//fmt.Println(coins)

	// don't need to run Tendermint updates if we exported
	if data.Exported {
		for _, lv := range data.LastValidatorPowers {
			keeper.SetLastValidatorPower(ctx, lv.Address, lv.Power)
			validator, found := keeper.GetValidator(ctx, lv.Address)
			if !found {
				panic(fmt.Sprintf("validator %s not found", lv.Address))
			}
			update := validator.ABCIValidatorUpdate()
			update.Power = lv.Power // keep the next-val-set offset, use the last power for the first block
			res = append(res, update)
		}
	} else {
		res = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	return res
}
