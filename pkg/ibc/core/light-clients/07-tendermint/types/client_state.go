package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
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

// Fraction defines the protobuf message types for tmmath.Fraction that only supports positive values.
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
// as a string array specifying the proof types for each position in chained proof
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

// Validate performs a basic validation of the client state fields.
func (cs ClientState) Validate() error {
	if strings.TrimSpace(cs.ChainId) == "" {
		return sdkerrors.Wrap(ErrInvalidChainID, "chain id cannot be empty string")
	}
	if cs.TrustingPeriod == 0 {
		return sdkerrors.Wrap(ErrInvalidTrustingPeriod, "trusting period cannot be zero")
	}
	if cs.UnbondingPeriod == 0 {
		return sdkerrors.Wrap(ErrInvalidUnbondingPeriod, "unbonding period cannot be zero")
	}
	if cs.MaxClockDrift == 0 {
		return sdkerrors.Wrap(ErrInvalidMaxClockDrift, "max clock drift cannot be zero")
	}
	if cs.LatestHeight.RevisionHeight == 0 {
		return sdkerrors.Wrapf(ErrInvalidHeaderHeight, "tendermint revision height cannot be zero")
	}
	if cs.TrustingPeriod >= cs.UnbondingPeriod {
		return sdkerrors.Wrapf(
			ErrInvalidTrustingPeriod,
			"trusting period (%s) should be < unbonding period (%s)", cs.TrustingPeriod, cs.UnbondingPeriod,
		)
	}

	if cs.ProofSpecs == nil {
		return sdkerrors.Wrap(ErrInvalidProofSpecs, "proof specs cannot be nil for tm client")
	}
	for i, spec := range cs.ProofSpecs {
		if spec == nil {
			return sdkerrors.Wrapf(ErrInvalidProofSpecs, "proof spec cannot be nil at index: %d", i)
		}
	}
	// UpgradePath may be empty, but if it isn't, each key must be non-empty
	for i, k := range cs.UpgradePath {
		if strings.TrimSpace(k) == "" {
			return sdkerrors.Wrapf(clienttypes.ErrInvalidClient, "key in upgrade path at index %d cannot be empty", i)
		}
	}

	return nil
}

func (cs ClientState) Initialize(ctx sdk.Context, clientStore sdk.KVStore, consState exported.ConsensusState) error {
	if _, ok := consState.(*ConsensusState); !ok {
		return sdkerrors.Wrapf(clienttypes.ErrInvalidConsensus, "invalid initial consensus state. expected type: %T, got: %T",
			&ConsensusState{}, consState)
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
		return sdkerrors.Wrap(clienttypes.ErrInvalidClient, "client state cannot be empty")
	}
	_, ok := clientState.(*ClientState)
	if !ok {
		return sdkerrors.Wrapf(clienttypes.ErrInvalidClient, "invalid client type %T, expected %T", clientState, &ClientState{})
	}

	bz, err := cdc.MarshalBinaryBare(clientState)
	if err != nil {
		return err
	}
	return merkleProof.VerifyMembership(cs.ProofSpecs, provingConsensusState.GetRoot(), path, bz)
}


// VerifyConnectionState verifies a proof of the connection state of the
// specified connection end stored on the target machine.
func (cs ClientState) VerifyConnectionState(
	store sdk.KVStore,
	cdc *codec.Codec,
	height exported.Height,
	prefix exported.Prefix,
	proof []byte,
	connectionID string,
	connectionEnd exported.ConnectionI,
) error {
	merkleProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}

	connectionPath := commitmenttypes.NewMerklePath(host.ConnectionPath(connectionID))
	path, err := commitmenttypes.ApplyPrefix(prefix, connectionPath)
	if err != nil {
		return err
	}

	connection, ok := connectionEnd.(connectiontypes.ConnectionEnd)
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid connection type %T", connectionEnd)
	}

	bz, err := cdc.MarshalBinaryBare(&connection)
	if err != nil {
		return err
	}

	if err := merkleProof.VerifyMembership(cs.ProofSpecs, consensusState.GetRoot(), path, bz); err != nil {
		return err
	}

	return nil
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
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "consensus state cannot be empty")
	}
	_, ok := consensusState.(*ConsensusState)
	if !ok {
		return sdkerrors.Wrapf(clienttypes.ErrInvalidConsensus, "invalid consensus type %T, expected %T", consensusState, &ConsensusState{})
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



func (cs ClientState) VerifyChannelState(store sdk.KVStore,
	cdc *codec.Codec, height exported.Height,
	prefix exported.Prefix, proof []byte,
	portID, channelID string, channel exported.ChannelI) error {
	merkleProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}
	channelPath := commitmenttypes.NewMerklePath(host.ChannelPath(portID, channelID))
	path, err := commitmenttypes.ApplyPrefix(prefix, channelPath)
	if err != nil {
		return err
	}
	channelEnd, ok := channel.(channeltypes.Channel)
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid channel type %T", channel)
	}

	bz, err := cdc.MarshalBinaryBare(&channelEnd)
	if err != nil {
		return err
	}

	if err := merkleProof.VerifyMembership(cs.ProofSpecs, consensusState.GetRoot(), path, bz); err != nil {
		return err
	}

	return nil
}

// todo ibc
func (cs ClientState) VerifyPacketAcknowledgement(store sdk.KVStore, cdc *codec.Codec, height exported.Height, currentTimestamp uint64, delayPeriod uint64, prefix exported.Prefix, proof []byte, portID, channelID string, sequence uint64, acknowledgement []byte) error {
	panic("implement me")
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
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d)", cs.GetLatestHeight(), height,
		)	}
	if cs.IsFrozen() && !cs.FrozenHeight.GT(height) {
		return commitmenttypes.MerkleProof{}, nil, clienttypes.ErrClientFrozen
	}
	if prefix == nil {
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrap(commitmenttypes.ErrInvalidPrefix, "prefix cannot be empty")
	}
	_, ok := prefix.(*commitmenttypes.MerklePrefix)
	if !ok {
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrapf(commitmenttypes.ErrInvalidPrefix, "invalid prefix type %T, expected *MerklePrefix", prefix)
	}

	if err = cdc.UnmarshalBinaryBare(proof, &merkleProof); err != nil {
		return commitmenttypes.MerkleProof{}, nil, sdkerrors.Wrap(commitmenttypes.ErrInvalidProof, "failed to unmarshal proof into commitment merkle proof")
	}
	consensusState, err = GetConsensusState(store, cdc, height)
	if err != nil {
		return commitmenttypes.MerkleProof{}, nil, err
	}

	return merkleProof, consensusState, nil
}


func (cs ClientState) VerifyPacketCommitment(store sdk.KVStore, cdc *codec.Codec, height exported.Height, currentTimestamp uint64, delayPeriod uint64, prefix exported.Prefix, proof []byte, portID, channelID string, sequence uint64, commitmentBytes []byte) error {
	merkleProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	if err := verifyDelayPeriodPassed(store, height, currentTimestamp, delayPeriod); err != nil {
		return err
	}

	commitmentPath := commitmenttypes.NewMerklePath(host.PacketCommitmentPath(portID, channelID, sequence))
	path, err := commitmenttypes.ApplyPrefix(prefix, commitmentPath)
	if err != nil {
		return err
	}

	if err := merkleProof.VerifyMembership(cs.ProofSpecs, consensusState.GetRoot(), path, commitmentBytes); err != nil {
		return err
	}

	return nil
}

func (cs ClientState) VerifyPacketReceiptAbsence(store sdk.KVStore, cdc *codec.Codec, height exported.Height, currentTimestamp uint64, delayPeriod uint64, prefix exported.Prefix, proof []byte, portID, channelID string, sequence uint64) error {
	merkleProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	if err := verifyDelayPeriodPassed(store, height, currentTimestamp, delayPeriod); err != nil {
		return err
	}

	receiptPath := commitmenttypes.NewMerklePath(host.PacketReceiptPath(portID, channelID, sequence))
	path, err := commitmenttypes.ApplyPrefix(prefix, receiptPath)
	if err != nil {
		return err
	}

	if err := merkleProof.VerifyNonMembership(cs.ProofSpecs, consensusState.GetRoot(), path); err != nil {
		return err
	}

	return nil
}

func (cs ClientState) VerifyNextSequenceRecv(store sdk.KVStore, cdc *codec.Codec, height exported.Height, currentTimestamp uint64, delayPeriod uint64, prefix exported.Prefix, proof []byte, portID, channelID string, nextSequenceRecv uint64) error {
	merkleProof, consensusState, err := produceVerificationArgs(store, cdc, cs, height, prefix, proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	if err := verifyDelayPeriodPassed(store, height, currentTimestamp, delayPeriod); err != nil {
		return err
	}

	nextSequenceRecvPath := commitmenttypes.NewMerklePath(host.NextSequenceRecvPath(portID, channelID))
	path, err := commitmenttypes.ApplyPrefix(prefix, nextSequenceRecvPath)
	if err != nil {
		return err
	}

	bz := sdk.Uint64ToBigEndian(nextSequenceRecv)

	if err := merkleProof.VerifyMembership(cs.ProofSpecs, consensusState.GetRoot(), path, bz); err != nil {
		return err
	}

	return nil
}


// verifyDelayPeriodPassed will ensure that at least delayPeriod amount of time has passed since consensus state was submitted
// before allowing verification to continue.
func verifyDelayPeriodPassed(store sdk.KVStore, proofHeight exported.Height, currentTimestamp, delayPeriod uint64) error {
	// check that executing chain's timestamp has passed consensusState's processed time + delay period
	processedTime, ok := GetProcessedTime(store, proofHeight)
	if !ok {
		return errors.Wrapf(ErrProcessedTimeNotFound, "processed time not found for height: %s", proofHeight)
	}
	validTime := processedTime + delayPeriod
	// NOTE: delay period is inclusive, so if currentTimestamp is validTime, then we return no error
	if validTime > currentTimestamp {
		return errors.Wrapf(ErrDelayPeriodNotPassed, "cannot verify packet until time: %d, current time: %d",
			validTime, currentTimestamp)
	}
	return nil
}

// GetProcessedTime gets the time (in nanoseconds) at which this chain received and processed a tendermint header.
// This is used to validate that a received packet has passed the delay period.
func GetProcessedTime(clientStore sdk.KVStore, height exported.Height) (uint64, bool) {
	key := ProcessedTimeKey(height)
	bz := clientStore.Get(key)
	if bz == nil {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz), true
}
