package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"strconv"

	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

// TODO-JT: carefully look at atomicity of this function
func (k Keeper) Attest(ctx sdk.Context, claim types.EthereumClaim, anyClaim *codectypes.Any) (*types.Attestation, error) {
	valAddr := k.GetOrchestratorValidator(ctx, claim.GetClaimer())
	if valAddr.Empty() {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well, but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce

	//lastEventNonce := k.GetLastEventNonceByValidator(ctx, valAddr)
	//if claim.GetEventNonce() != lastEventNonce+1 {
	//	return nil, types.ErrNonContiguousEventNonce
	//}

	k.Logger(ctx).Info("Attest", "EventNonce", claim.GetEventNonce(), "ClaimHash", hex.EncodeToString(claim.ClaimHash()))

	// Tries to get an attestation with the same eventNonce and claim as the claim that was submitted.
	att := k.GetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash())

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.Attestation{
			Observed: false,
			Height:   uint64(ctx.BlockHeight()),
			Claim:    anyClaim,
		}
	}

	// Add the validator's vote to this attestation
	att.Votes = append(att.Votes, valAddr.String())

	k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)
	k.setLastEventNonceByValidator(ctx, valAddr, claim.GetEventNonce())

	return att, nil
}

// TryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryAttestation(ctx sdk.Context, att *types.Attestation) {
	claim, err := k.UnpackAttestationClaim(att)
	if err != nil {
		panic("could not cast to claim")
	}
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.
	if !att.Observed {

		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
		requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
		k.Logger(ctx).Info("----- TryAttestation --- ", "totalPower", totalPower, "requiredPower", requiredPower)
		attestationPower := sdk.NewInt(0)
		for _, validator := range att.Votes {
			val:= sdk.HexToAddress(validator)
			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
			// Add it to the attestation power's sum
			attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))
			k.Logger(ctx).Info("----- Power Detail --- ", "validator", validator, "validatorPower", validatorPower, "AccumlateAttestationPower", attestationPower)
			// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
			// process the attestation, set Observed to true, and break
			if attestationPower.GTE(requiredPower) {
				lastEventNonce := k.GetLastObservedEventNonce(ctx)
				// this check is performed at the next level up so this should never panic
				// outside of programmer error.
				if claim.GetEventNonce() != uint64(lastEventNonce)+1 {
					panic("attempting to apply events to state out of order")
				}
				k.setLastObservedEventNonce(ctx, claim.GetEventNonce())
				k.SetLastObservedEthereumBlockHeight(ctx, claim.GetBlockHeight())

				att.Observed = true
				k.SetAttestation(ctx, claim.GetEventNonce(), claim.ClaimHash(), att)
				k.Logger(ctx).Info("----- ProcessAttestation --- ")
				k.processAttestation(ctx, att, claim)
				k.emitObservedEvent(ctx, att, claim)
				break
			}
		}
	} else {
		// We panic here because this should never happen
		panic("attempting to process observed attestation")
	}
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation, claim types.EthereumClaim) {
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, *att, claim); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", and validators can still be slashed for not
		// having voted for it.
		k.logger(ctx).Error("Attestation failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	} else {
		k.logger(ctx).Error("Attestation succeed")
		commit() // persist transient storage
	}
}

// emitObservedEvent emits an event with information about an attestation that has been applied to
// consensus state.
func (k Keeper) emitObservedEvent(ctx sdk.Context, att *types.Attestation, claim types.EthereumClaim) {
	observationEvent := sdk.NewEvent(
		types.EventTypeObservation,
		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
		sdk.NewAttribute([]byte(types.AttributeKeyAttestationType), []byte(string(claim.GetType()))),
		sdk.NewAttribute([]byte(types.AttributeKeyContract), []byte(k.GetBridgeContractAddress(ctx))),
		sdk.NewAttribute([]byte(types.AttributeKeyBridgeChainID), []byte(strconv.Itoa(int(k.GetBridgeChainID(ctx))))),
		sdk.NewAttribute([]byte(types.AttributeKeyAttestationID), []byte(string(types.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash())))), // todo: serialize with hex/ base64 ?
		sdk.NewAttribute([]byte(types.AttributeKeyNonce), []byte(fmt.Sprint(claim.GetEventNonce()))),
		// TODO: do we want to emit more information?
	)
	ctx.EventManager().EmitEvent(observationEvent)
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte, att *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	store.Set(aKey, codec.GetLegacyAminoByCodec(k.cdc).MustMarshalBinaryBare(att))
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte) *types.Attestation {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.Attestation
	codec.GetLegacyAminoByCodec(k.cdc).MustUnmarshalBinaryBare(bz, &att)
	return &att
}

// DeleteAttestation deletes an attestation given an event nonce and claim
func (k Keeper) DeleteAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte, att *types.Attestation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAttestationKeyWithHash(eventNonce, claimHash))
}

// GetAttestationMapping returns a mapping of eventnonce -> attestations at that nonce
func (k Keeper) GetAttestationMapping(ctx sdk.Context) (out map[uint64][]types.Attestation) {
	out = make(map[uint64][]types.Attestation)
	k.IterateAttestaions(ctx, func(_ []byte, att types.Attestation) bool {
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		if val, ok := out[claim.GetEventNonce()]; !ok {
			out[claim.GetEventNonce()] = []types.Attestation{att}
		} else {
			out[claim.GetEventNonce()] = append(val, att)
		}
		return false
	})
	return
}

// IterateAttestaions iterates through all attestations
func (k Keeper) IterateAttestaions(ctx sdk.Context, cb func([]byte, types.Attestation) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte(types.OracleAttestationKey)
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := types.Attestation{}
		codec.GetLegacyAminoByCodec(k.cdc).MustUnmarshalBinaryBare(iter.Value(), &att)
		// cb returns true to stop early
		if cb(iter.Key(), att) {
			return
		}
	}
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEventNonceKey)

	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

// GetLastObservedEthereumBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedEthereumBlockHeight(ctx sdk.Context) types.LastObservedEthereumBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEthereumBlockHeightKey)

	if len(bytes) == 0 {
		return types.LastObservedEthereumBlockHeight{
			CosmosBlockHeight:   0,
			EthereumBlockHeight: 0,
		}
	}
	height := types.LastObservedEthereumBlockHeight{}
	k.cdc.MustUnmarshalBinaryBare(bytes, &height)
	return height
}

// SetLastObservedEthereumBlockHeight sets the block height in the store.
func (k Keeper) SetLastObservedEthereumBlockHeight(ctx sdk.Context, ethereumHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	height := types.LastObservedEthereumBlockHeight{
		EthereumBlockHeight: ethereumHeight,
		CosmosBlockHeight:   uint64(ctx.BlockHeight()),
	}
	store.Set(types.LastObservedEthereumBlockHeightKey, k.cdc.MustMarshalBinaryBare(&height))
}

// setLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) setLastObservedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEventNonceKey, types.UInt64Bytes(nonce))
}

// GetLastEventNonceByValidator returns the latest event nonce for a given validator
func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByValidatorKey(validator))

	if len(bytes) == 0 {
		// in the case that we have no existing value this is the first
		// time a validator is submitting a claim. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		//
		// We could start at the LastObservedEventNonce but if we do that this
		// validator will be slashed, because they are responsible for making a claim
		// on any attestation that has not yet passed the slashing window.
		//
		// Therefore we need to return to them the lowest attestation that is still within
		// the slashing window. Since we delete attestations after the slashing window that's
		// just the lowest observed event in the store. If no claims have been submitted in for
		// params.SignedClaimsWindow we may have no attestations in our nonce. At which point
		// the last observed which is a persistant and never cleaned counter will suffice.
		lowest_observed := k.GetLastObservedEventNonce(ctx)
		attmap := k.GetAttestationMapping(ctx)
		// no new claims in params.SignedClaimsWindow, we can return the current value
		// because the validator can't be slashed for an event that has already passed.
		// so they only have to worry about the *next* event to occur
		if len(attmap) == 0 {
			return lowest_observed
		}
		for nonce, atts := range attmap {
			for att := range atts {
				if atts[att].Observed && nonce < lowest_observed {
					lowest_observed = nonce
				}
			}
		}
		// return the latest event minus one so that the validator
		// can submit that event and avoid slashing. special case
		// for zero
		if lowest_observed > 0 {
			return lowest_observed - 1
		} else {
			return 0
		}
	}
	return types.UInt64FromBytes(bytes)
}

// setLastEventNonceByValidator sets the latest event nonce for a give validator
func (k Keeper) setLastEventNonceByValidator(ctx sdk.Context, validator sdk.AccAddress, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByValidatorKey(validator), types.UInt64Bytes(nonce))
}
