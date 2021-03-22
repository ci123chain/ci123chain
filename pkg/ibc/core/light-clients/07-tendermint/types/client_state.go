package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	ics23 "github.com/confio/ics23/go"
	"github.com/pkg/errors"
	"strings"
	"time"
)

var _ exported.ClientState = (*ClientState)(nil)


// ClientState from Tendermint tracks the current validator set, latest height,
// and a possible frozen height.
type ClientState struct {
	ChainId    string   `json:"chain_id,omitempty"`
	TrustLevel Fraction `json:"trust_level" yaml:"trust_level"`
	// duration of the period since the LastestTimestamp during which the
	// submitted headers are valid for upgrade
	TrustingPeriod time.Duration `json:"trusting_period" yaml:"trusting_period"`
	// duration of the staking unbonding period
	UnbondingPeriod time.Duration `json:"unbonding_period" yaml:"unbonding_period"`
	// defines how much new (untrusted) header's Time can drift into the future.
	MaxClockDrift time.Duration `json:"max_clock_drift" yaml:"max_clock_drift"`
	// Block height when the client was frozen due to a misbehaviour
	FrozenHeight clienttypes.Height `json:"frozen_height" yaml:"frozen_height"`
	// Latest height the client was updated to
	LatestHeight clienttypes.Height `json:"latest_height" yaml:"latest_height"`
	// Proof specifications used in verifying counterparty state
	ProofSpecs []*ics23.ProofSpec `json:"proof_specs,omitempty" yaml:"proof_specs"`
	// Path at which next upgraded client will be committed.
	// Each element corresponds to the key for a single CommitmentProof in the chained proof.
	// NOTE: ClientState must stored under `{upgradePath}/{upgradeHeight}/clientState`
	// ConsensusState must be stored under `{upgradepath}/{upgradeHeight}/consensusState`
	// For SDK chains using the default upgrade module, upgrade_path should be []string{"upgrade", "upgradedIBCState"}`
	UpgradePath []string `json:"upgrade_path,omitempty" yaml:"upgrade_path"`
	// This flag, when set to true, will allow governance to recover a client
	// which has expired
	AllowUpdateAfterExpiry bool `json:"allow_update_after_expiry,omitempty" yaml:"allow_update_after_expiry"`
	// This flag, when set to true, will allow governance to unfreeze a client
	// whose chain has experienced a misbehaviour event
	AllowUpdateAfterMisbehaviour bool `json:"allow_update_after_misbehaviour,omitempty" yaml:"allow_update_after_misbehaviour"`
}


// Fraction defines the protobuf message type for tmmath.Fraction that only supports positive values.
type Fraction struct {
	Numerator   uint64 `json:"numerator,omitempty"`
	Denominator uint64 `json:"denominator,omitempty"`
}


// NewClientState creates a new ClientState instance
func NewClientState(
	chainID string, trustLevel Fraction,
	trustingPeriod, ubdPeriod, maxClockDrift time.Duration,
	latestHeight clienttypes.Height, specs []*ics23.ProofSpec,
	upgradePath []string, allowUpdateAfterExpiry, allowUpdateAfterMisbehaviour bool,
) *ClientState {
	return &ClientState{
		ChainId:                      chainID,
		TrustLevel:                   trustLevel,
		TrustingPeriod:               trustingPeriod,
		UnbondingPeriod:              ubdPeriod,
		MaxClockDrift:                maxClockDrift,
		LatestHeight:                 latestHeight,
		FrozenHeight:                 clienttypes.ZeroHeight(),
		ProofSpecs:                   specs,
		UpgradePath:                  upgradePath,
		AllowUpdateAfterExpiry:       allowUpdateAfterExpiry,
		AllowUpdateAfterMisbehaviour: allowUpdateAfterMisbehaviour,
	}
}


// GetChainID returns the chain-id
func (cs ClientState) GetChainID() string {
	return cs.ChainId
}

// ClientType is tendermint.
func (cs ClientState) ClientType() string {
	return exported.Tendermint
}

// GetLatestHeight returns latest block height.
func (cs ClientState) GetLatestHeight() exported.Height {
	return cs.LatestHeight
}

// IsFrozen returns true if the frozen height has been set.
func (cs ClientState) IsFrozen() bool {
	return !cs.FrozenHeight.IsZero()
}


// GetFrozenHeight returns the height at which client is frozen
// NOTE: FrozenHeight is zero if client is unfrozen
func (cs ClientState) GetFrozenHeight() exported.Height {
	return cs.FrozenHeight
}

// IsExpired returns whether or not the client has passed the trusting period since the last
// update (in which case no headers are considered valid).
func (cs ClientState) IsExpired(latestTimestamp, now time.Time) bool {
	expirationTime := latestTimestamp.Add(cs.TrustingPeriod)
	return !expirationTime.After(now)
}

// GetProofSpecs returns the format the client expects for proof verification
// as a string array specifying the proof type for each position in chained proof
func (cs ClientState) GetProofSpecs() []*ics23.ProofSpec {
	return cs.ProofSpecs
}

// ZeroCustomFields returns a ClientState that is a copy of the current ClientState
// with all client customizable fields zeroed out
func (cs ClientState) ZeroCustomFields() exported.ClientState {
	// copy over all chain-specified fields
	// and leave custom fields empty
	return &ClientState{
		ChainId:         cs.ChainId,
		UnbondingPeriod: cs.UnbondingPeriod,
		LatestHeight:    cs.LatestHeight,
		ProofSpecs:      cs.ProofSpecs,
		UpgradePath:     cs.UpgradePath,
	}
}

// basic validate for cs
func (cs ClientState) Validate() error {
	if strings.TrimSpace(cs.ChainId) == "" {
		return errors.New("chain-id cannot be empty string")
	}
	// todo verify tendermint trustlevel

	if cs.TrustingPeriod == 0 {
		return errors.New("trusting period cannot be zero")
	}
	if cs.UnbondingPeriod == 0 {
		return errors.New("unbonding period cannot be zero")
	}
	if cs.MaxClockDrift == 0 {
		return errors.New("max clock drift cannot be zero")
	}
	if cs.LatestHeight.RevisionHeight == 0 {
		return errors.New("tendermint revision height cannot be zero")
	}
	if cs.TrustingPeriod >= cs.UnbondingPeriod {
		return fmt.Errorf("trusting period (%s) should be < unbounding period (%s)", cs.TrustingPeriod, cs.UnbondingPeriod)
	}

	if cs.ProofSpecs == nil {
		return errors.New("proof specs cannot be nil")
	}

	for i, spec := range cs.ProofSpecs {
		if  spec == nil {
			return fmt.Errorf("proof spec cannot be nil at index %d", i)
		}
	}

	for i, k := range cs.UpgradePath {
		if strings.TrimSpace(k) == "" {
			return errors.Errorf("key in upgrade path at index %d cannot be empty string", i)
		}
	}
	return nil
}

func (cs ClientState) Initialize(ctx sdk.Context, clientStore sdk.KVStore, consState exported.ConsensusState) error {
	if _, ok := consState.(*ConsensusState); !ok {
		return errors.Errorf("invalid initial consensus state. expected type: %T, got: %T", &ConsensusState{}, consState)
	}
	SetProcessedTime(clientStore, cs.GetLatestHeight(), uint64(ctx.BlockHeader().Time.UnixNano()))
	return nil
}

func (cs ClientState) CheckHeaderAndUpdateState(ctx sdk.Context, cdc *codec.Codec, store sdk.KVStore,
	header exported.Header) (exported.ClientState, exported.ConsensusState, error) {
	return nil, nil, nil
}


func (cs ClientState) VerifyUpgradeAndUpdateState(ctx sdk.Context,
	cdc *codec.Codec, store sdk.KVStore,
	newClient exported.ClientState, proofUpgradeClient,
	proofUpgradeConsState []byte) (exported.ClientState, exported.ConsensusState, error) {
	return nil, nil, nil
}


func (cs ClientState) VerifyClientState(
	store sdk.KVStore,
	cdc *codec.Codec,
	height exported.Height,
	prefix exported.Prefix,
	counterpartyClientIdentifier string,
	proof []byte,
	clientState exported.ClientState) error {

	merkleProof, provingConsensusState, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}
	clientPrefixedPath := commitmenttypes.NewMerklePath(host.FullClientStatePath(counterpartyClientIdentifier))
	path, err := commitmenttypes.ApplyPrefix(prefix, clientPrefixedPath)
	if err != nil {
		return err
	}
	if clientState == nil {
		return errors.New("client state cannot be empty")
	}
	_, ok := clientState.(*ClientState)
	if !ok {
		return errors.Errorf("invalid client type %T, expected %T", clientState, &ClientState{})
	}

	bz, err := cdc.MarshalBinaryBare(clientState)
	if err != nil {
		return err
	}
	return merkleProof.VerifyMembership(cs.ProofSpecs, provingConsensusState.GetRoot(), path, bz)
}

func (cs ClientState) VerifyClientConsensusState(
	store sdk.KVStore,
	cdc *codec.Codec,
	height exported.Height,
	counterpartyClientIdentifier string,
	consensusHeight exported.Height, // todo ?
	prefix exported.Prefix,
	proof []byte,
	consensusState exported.ConsensusState,
) error {
	merkleProof, provingConsensusState, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}
	clientPrefixedPath := commitmenttypes.NewMerklePath(host.FullConsensusStatePath(counterpartyClientIdentifier, consensusHeight))
	path, err := commitmenttypes.ApplyPrefix(prefix, clientPrefixedPath)
	if err != nil {
		return err
	}
	if consensusState == nil {
		return errors.New("consensus state cannot be empty")
	}
	_, ok := consensusState.(*ConsensusState)
	if !ok {
		return errors.Errorf("invalid consensus type %T, expected %T", consensusState, &ConsensusState{})
	}
	bz, err := cdc.MarshalBinaryBare(consensusState)
	if err != nil {
		return err
	}
	if err := merkleProof.VerifyMembership(cs.ProofSpecs, provingConsensusState.GetRoot(), path, bz); err != nil {
		return err
	}
	return nil
}



// private methods implements

// produceVerificationArgs perfoms the basic checks on the arguments that are
// shared between the verification functions and returns the unmarshalled
// merkle proof, the consensus state and an error if one occurred.
func produceVerificationArgs(
	store sdk.KVStore,
	cdc *codec.Codec,
	cs ClientState,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
) (merkleProof commitmenttypes.MerkleProof, consensusState *ConsensusState, err error) {
	if cs.GetLatestHeight().LT(height) {
		return commitmenttypes.MerkleProof{}, nil, errors.Errorf("client state height < proof height (%d < %d)", cs.GetLatestHeight(), height)
	}
	if cs.IsFrozen() && !cs.FrozenHeight.GT(height) {
		return commitmenttypes.MerkleProof{}, nil, errors.New("light client is frozen due to misbehaviou")
	}
	if prefix == nil {
		return commitmenttypes.MerkleProof{}, nil, errors.New("prefix cannot be empty")
	}
	_, ok := prefix.(*commitmenttypes.MerklePrefix)
	if !ok {
		return commitmenttypes.MerkleProof{}, nil, errors.Errorf("invalid prefix type %T, expected *MerklePrefix", prefix)
	}

	if err = cdc.UnmarshalBinaryBare(proof, &merkleProof); err != nil {
		return commitmenttypes.MerkleProof{}, nil, errors.New("failed to unmarshal proof into commitment merkle proof")
	}
	consensusState, err = GetConsensusState(store, cdc, height)
	if err != nil {
		return commitmenttypes.MerkleProof{}, nil, err
	}

	return merkleProof, consensusState, nil
}