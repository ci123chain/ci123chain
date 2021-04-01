package exported

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	ics23 "github.com/confio/ics23/go"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

const (
	// Solomachine is used to indicate that the light client is a solo machine.
	Solomachine string = "06-solomachine"

	// Tendermint is used to indicate that the client uses the Tendermint Consensus Algorithm.
	Tendermint string = "07-tendermint"

	// Localhost is the client types for a localhost client. It is also used as the clientID
	Localhost string = "09-localhost"
)

type ClientState interface {
	ClientType() string
	GetLatestHeight() Height
	IsFrozen() bool
	GetFrozenHeight() Height
	Validate() error
	GetProofSpecs() []*ics23.ProofSpec

	Initialize(ctx sdk.Context, store sdk.KVStore, cs ConsensusState) error

	CheckHeaderAndUpdateState(sdk.Context, *codec.Codec, sdk.KVStore, Header) (ClientState, ConsensusState, error)
	VerifyUpgradeAndUpdateState(
		ctx sdk.Context,
		cdc *codec.Codec,
		store sdk.KVStore,
		newClient ClientState,
		proofUpgradeClient,
		proofUpgradeConsState []byte,
	) (ClientState, ConsensusState, error)

	// 清除用户自定义字段
	ZeroCustomFields() ClientState

	VerifyClientState(
		store sdk.KVStore,
		cdc *codec.Codec,
		height Height,
		prefix Prefix,
		counterpartyClientIdentifier string,
		proof []byte,
		clientState ClientState,
	) error

	VerifyClientConsensusState(
		store sdk.KVStore,
		cdc *codec.Codec,
		height Height,
		counterpartyClientIdentifier string,
		consensusHeight Height,
		prefix Prefix,
		proof []byte,
		consensusState ConsensusState,
	) error

	VerifyConnectionState(
		store sdk.KVStore,
		cdc *codec.Codec,
		height Height,
		prefix Prefix,
		proof []byte,
		connectionID string,
		connectionEnd ConnectionI,
	) error

	VerifyChannelState(
		store sdk.KVStore,
		cdc *codec.Codec,
		height Height,
		prefix Prefix,
		proof []byte,
		portID,
		channelID string,
		channel ChannelI,
	) error

	VerifyPacketAcknowledgement(
		store sdk.KVStore,
		cdc *codec.Codec,
		height Height,
		currentTimestamp uint64,
		delayPeriod uint64,
		prefix Prefix,
		proof []byte,
		portID,
		channelID string,
		sequence uint64,
		acknowledgement []byte,
	) error
}

type ConsensusState interface {
	ClientType() string
	GetRoot() Root
	GetTimestamp() uint64
	ValidateBasic() error
}

// Header is the consensus state update information
type Header interface {
	ClientType() string
	GetHeight() Height
	ValidateBasic() error
}


// Height is a wrapper interface over clienttypes.Height
// all clients must use the concrete implementation in types
type Height interface {
	IsZero() bool
	LT(Height) bool
	LTE(Height) bool
	EQ(Height) bool
	GT(Height) bool
	GTE(Height) bool
	GetRevisionNumber() uint64
	GetRevisionHeight() uint64
	Increment() Height
	Decrement() (Height, bool)
	String() string
}
