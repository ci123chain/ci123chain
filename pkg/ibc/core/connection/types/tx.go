package types

import (
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)


type MsgConnectionOpenInit struct {
	ClientId     string       `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty" yaml:"client_id"`
	Counterparty Counterparty `protobuf:"bytes,2,opt,name=counterparty,proto3" json:"counterparty"`
	Version      *Version     `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	DelayPeriod  uint64       `protobuf:"varint,4,opt,name=delay_period,json=delayPeriod,proto3" json:"delay_period,omitempty" yaml:"delay_period"`
	Signer       string       `protobuf:"bytes,5,opt,name=signer,proto3" json:"signer,omitempty"`
}



type MsgConnectionOpenInitResponse struct {
}


// MsgConnectionOpenTry defines a msg sent by a Relayer to try to open a
// connection on Chain B.
type MsgConnectionOpenTry struct {
	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty" yaml:"client_id"`
	// in the case of crossing hello's, when both chains call OpenInit, we need the connection identifier
	// of the previous connection in state INIT
	PreviousConnectionId string        `protobuf:"bytes,2,opt,name=previous_connection_id,json=previousConnectionId,proto3" json:"previous_connection_id,omitempty" yaml:"previous_connection_id"`
	ClientState          exported.ClientState    `protobuf:"bytes,3,opt,name=client_state,json=clientState,proto3" json:"client_state,omitempty" yaml:"client_state"`
	Counterparty         Counterparty  `protobuf:"bytes,4,opt,name=counterparty,proto3" json:"counterparty"`
	DelayPeriod          uint64        `protobuf:"varint,5,opt,name=delay_period,json=delayPeriod,proto3" json:"delay_period,omitempty" yaml:"delay_period"`
	CounterpartyVersions []*Version    `protobuf:"bytes,6,rep,name=counterparty_versions,json=counterpartyVersions,proto3" json:"counterparty_versions,omitempty" yaml:"counterparty_versions"`
	ProofHeight          clienttypes.Height `protobuf:"bytes,7,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
	// proof of the initialization the connection on Chain A: `UNITIALIZED ->
	// INIT`
	ProofInit []byte `protobuf:"bytes,8,opt,name=proof_init,json=proofInit,proto3" json:"proof_init,omitempty" yaml:"proof_init"`
	// proof of client state included in message
	ProofClient []byte `protobuf:"bytes,9,opt,name=proof_client,json=proofClient,proto3" json:"proof_client,omitempty" yaml:"proof_client"`
	// proof of client consensus state
	ProofConsensus  []byte        `protobuf:"bytes,10,opt,name=proof_consensus,json=proofConsensus,proto3" json:"proof_consensus,omitempty" yaml:"proof_consensus"`
	ConsensusHeight clienttypes.Height `protobuf:"bytes,11,opt,name=consensus_height,json=consensusHeight,proto3" json:"consensus_height" yaml:"consensus_height"`
	Signer          string        `protobuf:"bytes,12,opt,name=signer,proto3" json:"signer,omitempty"`
}


// MsgConnectionOpenTryResponse defines the Msg/ConnectionOpenTry response types.
type MsgConnectionOpenTryResponse struct {
}



// MsgConnectionOpenAck defines a msg sent by a Relayer to Chain A to
// acknowledge the change of connection state to TRYOPEN on Chain B.
type MsgConnectionOpenAck struct {
	ConnectionId             string        `protobuf:"bytes,1,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty" yaml:"connection_id"`
	CounterpartyConnectionId string        `protobuf:"bytes,2,opt,name=counterparty_connection_id,json=counterpartyConnectionId,proto3" json:"counterparty_connection_id,omitempty" yaml:"counterparty_connection_id"`
	Version                  *Version      `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	ClientState              exported.ClientState    `protobuf:"bytes,4,opt,name=client_state,json=clientState,proto3" json:"client_state,omitempty" yaml:"client_state"`
	ProofHeight              clienttypes.Height `protobuf:"bytes,5,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
	// proof of the initialization the connection on Chain B: `UNITIALIZED ->
	// TRYOPEN`
	ProofTry []byte `protobuf:"bytes,6,opt,name=proof_try,json=proofTry,proto3" json:"proof_try,omitempty" yaml:"proof_try"`
	// proof of client state included in message
	ProofClient []byte `protobuf:"bytes,7,opt,name=proof_client,json=proofClient,proto3" json:"proof_client,omitempty" yaml:"proof_client"`
	// proof of client consensus state
	ProofConsensus  []byte        `protobuf:"bytes,8,opt,name=proof_consensus,json=proofConsensus,proto3" json:"proof_consensus,omitempty" yaml:"proof_consensus"`
	ConsensusHeight clienttypes.Height `protobuf:"bytes,9,opt,name=consensus_height,json=consensusHeight,proto3" json:"consensus_height" yaml:"consensus_height"`
	Signer          string        `protobuf:"bytes,10,opt,name=signer,proto3" json:"signer,omitempty"`
}


// MsgConnectionOpenAckResponse defines the Msg/ConnectionOpenAck response types.
type MsgConnectionOpenAckResponse struct {
}


// MsgConnectionOpenConfirm defines a msg sent by a Relayer to Chain B to
// acknowledge the change of connection state to OPEN on Chain A.
type MsgConnectionOpenConfirm struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty" yaml:"connection_id"`
	// proof for the change of the connection state on Chain A: `INIT -> OPEN`
	ProofAck    []byte        `protobuf:"bytes,2,opt,name=proof_ack,json=proofAck,proto3" json:"proof_ack,omitempty" yaml:"proof_ack"`
	ProofHeight clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height" yaml:"proof_height"`
	Signer      string        `protobuf:"bytes,4,opt,name=signer,proto3" json:"signer,omitempty"`
}

// MsgConnectionOpenConfirmResponse defines the Msg/ConnectionOpenConfirm response types.
type MsgConnectionOpenConfirmResponse struct {
}
