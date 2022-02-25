package gravity

import (
	"encoding/json"
	"sort"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/gravity/keeper"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Question: what here can be epoched?
	slashing(ctx, k)
	attestationTally(ctx, k)
	cleanupTimedOutBatches(ctx, k)
	cleanupTimedOutLogicCalls(ctx, k)
	cleanupConfirmedValsets(ctx, k)
}

func cleanupConfirmedValsets(ctx sdk.Context, k keeper.Keeper) {
	// Auto signerset tx creation.
	// 1. If there are no signer set requests, create a new one.
	// 2. If there is at least one validator who started unbonding in current block. (we persist last unbonded block height in hooks.go)
	//      This will make sure the unbonding validator has to provide an ethereum signature to a new signer set tx
	//	    that excludes him before he completely Unbonds.  Otherwise he will be slashed
	// 3. If power change between validators of Current signer set and latest signer set request is > 5%
	valsets := k.GetValsets(ctx)
	if len(valsets) == 0 {
		k.SetValsetRequest(ctx)
		return
	}

	latestConfirmNonce := k.GetLastValsetConfirmNonce(ctx)
	if latestConfirmNonce == 0 {
		return
	}

	lastUnbondingHeight := k.GetLastUnBondingBlockHeight(ctx)
	blockHeight := uint64(ctx.BlockHeight())
	
	powerDiff := types.BridgeValidators(k.GetCurrentValset(ctx).Members).PowerDiff(valsets[0].Members)

	shouldCreate := (lastUnbondingHeight == blockHeight) || (powerDiff > 0.05)
	k.Logger(ctx).Info(
		"considering signer set tx creation",
		"blockHeight", blockHeight,
		"lastUnbondingHeight", lastUnbondingHeight,
		"latestSignerSetTx.Nonce", latestConfirmNonce,
		"powerDiff", powerDiff,
		"shouldCreate", shouldCreate,
	)

	if shouldCreate {
		k.SetValsetRequest(ctx)
	}

	k.DeleteValset(ctx, latestConfirmNonce-1)
}

func slashing(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)
	currentBondedSet := k.StakingKeeper.GetBondedValidatorsByPower(ctx)

	//valsets are sorted so the most recent one is first
	//valsets := k.GetValsets(ctx)
	//if len(valsets) == 0 {
	//	k.SetValsetRequest(ctx)
	//}
	//
	//for i, vs := range valsets {
	//	signedWithinWindow := uint64(ctx.BlockHeight()) > params.SignedValsetsWindow && uint64(ctx.BlockHeight())-params.SignedValsetsWindow > vs.Height
	//	switch {
	//	// #1 condition
	//	// We look through the full bonded validator set (not just the active set, include unbonding validators)
	//	// and we slash users who haven't signed a valset that is currentHeight - signedBlocksWindow old
	//	case signedWithinWindow:
	//		// keep the latest valset
	//		// first we need to see which validators in the active set // haven't signed the valdiator set and slash them,
	//		confirms := k.GetValsetConfirms(ctx, vs.Nonce)
	//
	//		// do nothing for only one validator
	//		for _, val := range currentBondedSet {
	//			found := false
	//			for _, conf := range confirms {
	//				if conf.EthAddress == k.GetEthAddress(ctx, val.GetOperator()) {
	//					found = true
	//					break
	//				}
	//			}
	//			if !found {
	//				cons := val.GetConsAddr()
	//				k.StakingKeeper.Slash(ctx, cons,
	//					ctx.BlockHeight(), val.ConsensusPower(),
	//					params.SlashFractionValset)
	//				//k.StakingKeeper.Jail(ctx, cons)
	//			}
	//		}
	//		// then we prune the valset from state
	//		k.DeleteValset(ctx, vs.Nonce)
	//	}
	//	// on the latest validator set, check for change in power against
	//	// current, and emit a new validator set if the change in power >5%
	//	//case i == 0:
	//	if i == 0 && types.BridgeValidators(k.GetCurrentValset(ctx).Members).PowerDiff(vs.Members) > 0.05 {
	//		k.SetValsetRequest(ctx)
	//	}
	//}


	// #2 condition
	// We look through the full bonded set (not just the active set, include unbonding validators)
	// and we slash users who haven't signed a batch confirmation that is >15hrs in blocks old
	batches := k.GetOutgoingTxBatches(ctx)
	for _, batch := range batches {
		signedWithinWindow := uint64(ctx.BlockHeight()) > params.SignedBatchesWindow && uint64(ctx.BlockHeight())-params.SignedBatchesWindow > batch.Block
		if signedWithinWindow {
			//confirms := k.GetBatchConfirmByNonceAndTokenContract(ctx, batch.BatchNonce, batch.TokenContract)
			//for _, val := range currentBondedSet {
			//	found := false
			//	for _, conf := range confirms {
			//		// TODO: double check this logic
			//		confVal, _ := sdk.AccAddressFromBech32(conf.Orchestrator)
			//		if k.GetOrchestratorValidator(ctx, confVal).Equals(val.GetOperator()) {
			//			found = true
			//			break
			//		}
			//	}
			//	if !found {
			//		cons := val.GetConsAddr()
			//		k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), val.ConsensusPower(), params.SlashFractionBatch)
			//		k.StakingKeeper.Jail(ctx, cons)
			//	}
			//}

			// clean up batches here
			k.DeleteBatch(ctx, *batch)
		}
	}

	// #3 condition
	// Oracle events MsgDepositClaim, MsgWithdrawClaim
	attmap := k.GetAttestationMapping(ctx)
	for _, atts := range attmap {
		// slash conflicting votes
		if len(atts) > 1 {
			var unObs []types.Attestation
			oneObserved := false
			for _, att := range atts {
				if att.Observed {
					oneObserved = true
					continue
				}
				unObs = append(unObs, att)
			}
			// if one is observed delete the *other* attestations, do not delete
			// the original as we will need it later.
			if oneObserved {
				for _, att := range unObs {
					//for _, valaddr := range att.Votes {
					//	validator := sdk.HexToAddress(valaddr)
					//	val := k.StakingKeeper.Validator(ctx, validator)
					//	cons := val.GetConsAddr()
					//	k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), k.StakingKeeper.GetLastValidatorPower(ctx, validator), params.SlashFractionConflictingClaim)
					//	k.StakingKeeper.Jail(ctx, cons)
					//}
					claim, err := k.UnpackAttestationClaim(&att)
					if err != nil {
						panic("couldn't cast to claim")
					}

					k.DeleteAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
				}
			}
		}

		if len(atts) == 1 {
			att := atts[0]
			// TODO-JT: Review this
			windowPassed := uint64(ctx.BlockHeight()) > params.SignedClaimsWindow &&
				uint64(ctx.BlockHeight())-params.SignedClaimsWindow > att.Height

			// if the signing window has passed and the attestation is still unobserved wait.
			if windowPassed && att.Observed {
				for _, bv := range currentBondedSet {
					found := false
					for _, val := range att.Votes {
						confVal := sdk.HexToAddress(val)
						if confVal.Equals(bv.GetOperator()) {
							found = true
							break
						}
					}
					if !found {
						cons := bv.GetConsAddr()
						k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), k.StakingKeeper.GetLastValidatorPower(ctx, bv.GetOperator()), params.SlashFractionClaim)
						k.StakingKeeper.Jail(ctx, cons)
					}
				}
				claim, err := k.UnpackAttestationClaim(&att)
				if err != nil {
					panic("couldn't cast to claim")
				}

				k.DeleteAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), &att)
			}
		}
	}
}


// Iterate over all attestations currently being voted on in order of nonce and
// "Observe" those who have passed the threshold. Break the loop once we see
// an attestation that has not passed the threshold
func attestationTally(ctx sdk.Context, k keeper.Keeper) {
	attmap := k.GetAttestationMapping(ctx)
	// We make a slice with all the event nonces that are in the attestation mapping
	keys := make([]uint64, 0, len(attmap))
	for k := range attmap {
		keys = append(keys, k)
	}
	// Then we sort it
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	// This iterates over all keys (event nonces) in the attestation mapping. Each value contains
	// a slice with one or more attestations at that event nonce. There can be multiple attestations
	// at one event nonce when validators disagree about what event happened at that nonce.
	for _, nonce := range keys {
		// This iterates over all attestations at a particular event nonce.
		// They are ordered by when the first attestation at the event nonce was received.
		// This order is not important.
		for _, att := range attmap[nonce] {
			// We check if the event nonce is exactly 1 higher than the last attestation that was
			// observed. If it is not, we just move on to the next nonce. This will skip over all
			// attestations that have already been observed.
			//
			// Once we hit an event nonce that is one higher than the last observed event, we stop
			// skipping over this conditional and start calling tryAttestation (counting votes)
			// Once an attestation at a given event nonce has enough votes and becomes observed,
			// every other attestation at that nonce will be skipped, since the lastObservedEventNonce
			// will be incremented.
			//
			// Then we go to the next event nonce in the attestation mapping, if there is one. This
			// nonce will once again be one higher than the lastObservedEventNonce.
			// If there is an attestation at this event nonce which has enough votes to be observed,
			// we skip the other attestations and move on to the next nonce again.
			// If no attestation becomes observed, when we get to the next nonce, every attestation in
			// it will be skipped. The same will happen for every nonce after that.
			lastEventNonce := k.GetLastObservedEventNonce(ctx)
			claim, _ := k.UnpackAttestationClaim(&att)
			claimbz, _ := json.Marshal(claim)

			if nonce == uint64(lastEventNonce)+1 {
				k.Logger(ctx).Info("TryAttestation", "Nonce", nonce, "lastEventNonce", lastEventNonce, "ClaimJson", string(claimbz))
				k.TryAttestation(ctx, &att)
			} else {
				k.Logger(ctx).Info("Try Not Attestation", "Nonce", nonce, "lastEventNonce", lastEventNonce, "ClaimJson", string(claimbz))
			}
		}
	}
}

// cleanupTimedOutBatches deletes batches that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning batch 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single batch at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutBatches(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx).EthereumBlockHeight
	batches := k.GetOutgoingTxBatches(ctx)
	for _, batch := range batches {
		if batch.BatchTimeout < ethereumHeight {
			k.CancelOutgoingTXBatch(ctx, batch.TokenContract, batch.BatchNonce)
		}
	}
}

// cleanupTimedOutBatches deletes logic calls that have passed their expiration on Ethereum
// keep in mind several things when modifying this function
// A) unlike nonces timeouts are not monotonically increasing, meaning call 5 can have a later timeout than batch 6
//    this means that we MUST only cleanup a single call at a time
// B) it is possible for ethereumHeight to be zero if no events have ever occurred, make sure your code accounts for this
// C) When we compute the timeout we do our best to estimate the Ethereum block height at that very second. But what we work with
//    here is the Ethereum block height at the time of the last Deposit or Withdraw to be observed. It's very important we do not
//    project, if we do a slowdown on ethereum could cause a double spend. Instead timeouts will *only* occur after the timeout period
//    AND any deposit or withdraw has occurred to update the Ethereum block height.
func cleanupTimedOutLogicCalls(ctx sdk.Context, k keeper.Keeper) {
	ethereumHeight := k.GetLastObservedEthereumBlockHeight(ctx).EthereumBlockHeight
	calls := k.GetOutgoingLogicCalls(ctx)
	for _, call := range calls {
		if call.Timeout < ethereumHeight {
			k.CancelOutgoingLogicCall(ctx, call.InvalidationId, call.InvalidationNonce)
		}
	}
}

// TestingEndBlocker is a second endblocker function only imported in the Gravity codebase itself
// if you are a consuming Cosmos chain DO NOT IMPORT THIS, it simulates a chain using the arbitrary
// logic API to request logic calls
func TestingEndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// if this is nil we have not set our test outgoing logic call yet
	if k.GetOutgoingLogicCall(ctx, []byte("GravityTesting"), 0).Payload == nil {
		// TODO this call isn't actually very useful for testing, since it always
		// throws, being just junk data that's expected. But it prevents us from checking
		// the full lifecycle of the call. We need to find some way for this to read data
		// and encode a simple testing call, probably to one of the already deployed ERC20
		// contracts so that we can get the full lifecycle.
		token := []*types.ERC20Token{{
			Contract: "0x7580bfe88dd3d07947908fae12d95872a260f2d8",
			Amount:   sdk.NewIntFromUint64(5000),
		}}
		_ = types.OutgoingLogicCall{
			Transfers:            token,
			Fees:                 token,
			LogicContractAddress: "0x510ab76899430424d209a6c9a5b9951fb8a6f47d",
			Payload:              []byte("fake bytes"),
			Timeout:              10000,
			InvalidationId:       []byte("GravityTesting"),
			InvalidationNonce:    1,
		}
		//k.SetOutgoingLogicCall(ctx, &call)
	}
}
