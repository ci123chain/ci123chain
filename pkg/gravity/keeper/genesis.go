package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, *data.Params)
	// reset valsets in state
	for _, vs := range data.Valsets {
		// TODO: block height?
		k.StoreValsetUnsafe(ctx, vs)
	}

	for gid, data := range data.Gravitys {
		k.SetCurrentGid(gid)
		k.saveGravityID(ctx, gid)
		// reset valset confirmations in state
		for _, conf := range data.ValsetConfirms {
			k.SetValsetConfirmByGID(ctx, *conf)
		}

		// reset batches in state
		for _, batch := range data.Batches {
			// TODO: block height?
			k.StoreBatchUnsafe(ctx, batch)
		}

		// reset batch confirmations in state
		for _, conf := range data.BatchConfirms {
			k.SetBatchConfirmWithGID(ctx, &conf)
		}


		// reset pool transactions in state
		for _, tx := range data.UnbatchedTransfers {
			k.setPoolEntry(ctx, tx)
			k.setTxIdState(ctx, tx.Id, txIdStatePending)
		}

		// reset attestations in state
		for _, att := range data.Attestations {
			claim, err := k.UnpackAttestationClaim(&att)
			if err != nil {
				panic("couldn't cast to claim")
			}

			// TODO: block height?
			k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
		}
		k.setLastObservedEventNonceWithGid(ctx, data.LastObservedNonce)

		// reset attestation state of specific validators
		// this must be done after the above to be correct
		for _, att := range data.Attestations {
			claim, err := k.UnpackAttestationClaim(&att)
			if err != nil {
				panic("couldn't cast to claim")
			}

			for _, vote := range att.Votes {
				val:= sdk.HexToAddress(vote)
				last := k.GetLastEventNonceByValidator(ctx, val)
				if claim.GetEventNonce() > last {
					k.setLastEventNonceByValidator(ctx, val, claim.GetEventNonce())
				}
			}
		}

		// populate state with cosmos originated denom-erc20 mapping
		for _, item := range data.Erc20ToDenoms {
			k.setERC20Map(ctx, item.Denom, item.Erc20)
		}

		k.SetCurrentGid("")
	}


}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var (
		p                   = k.GetParams(ctx)
		//calls               = k.GetOutgoingLogicCalls(ctx)

		valsets             = k.GetValsets(ctx)
		vsconfs             = []*types.MsgValsetConfirm{}
		batchconfs          = []types.MsgConfirmBatch{}
		attestations        = []types.Attestation{}
		delegates           = k.GetDelegateKeys(ctx)

		erc20ToDenoms       = []*types.ERC20ToDenom{}

		gravityDatas		= map[string]types.GravityData{}
	)
	gids := k.GetAllGravityIDs(ctx)
	for _, gid := range gids {
		k.SetCurrentGid(gid)

		lastobserved := k.GetLastObservedEventNonceWithGid(ctx)

		//export valset confirmations from state
		for _, vs := range valsets {
			// TODO: set height = 0?
			vsconfs = append(vsconfs, k.GetValsetConfirmsByGID(ctx, vs.Nonce)...)
		}

		// export batch confirmations from state
		batches := k.GetOutgoingTxBatches(ctx)
		for _, batch := range batches {
			// TODO: set height = 0?
			batchconfs = append(batchconfs, k.GetBatchConfirmByNonceAndTokenContractWithGID(ctx, batch.BatchNonce, batch.TokenContract)...)
		}

		attmap := k.GetAttestationMapping(ctx)
		// export attestations from state
		for _, atts := range attmap {
			// TODO: set height = 0?
			attestations = append(attestations, atts...)
		}

		// export erc20 to denom relations
		k.IterateERC20ToDenom(ctx, func(key []byte, erc20ToDenom *types.ERC20ToDenom) bool {
			erc20ToDenoms = append(erc20ToDenoms, erc20ToDenom)
			return false
		})

		unbatched_transfers := k.GetPoolTransactions(ctx)

		gData := types.GravityData{
			LastObservedNonce:  lastobserved,
			ValsetConfirms:     vsconfs,
			Batches:            batches,
			BatchConfirms:      batchconfs,
			Attestations:       attestations,
			Erc20ToDenoms:      erc20ToDenoms,
			UnbatchedTransfers: unbatched_transfers,
		}
		gravityDatas[gid] = gData

		k.SetCurrentGid("")
	}

	return types.GenesisState{
		Params:             &p,
		Valsets:            valsets,
		DelegateKeys:       delegates,
		Gravitys: 			gravityDatas,
	}
}
