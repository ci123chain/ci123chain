package types

import (
	"errors"
	abcitypes "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
)

const (
	TypeMsgCreateClient		string = "create_client"
)

type MsgCreateClient struct {
	// client_state
	ClientState exported.ClientState `json:"client_state,omitempty" yaml:"client_state"`
	// consensus_state
	ConsensusState exported.ConsensusState `json:"consensus_state,omitempty" yaml:"consensus_state"`
	// singer address
	Signer               string   `json:"signer,omitempty"`
}

func NewMsgCreateClient() (*MsgCreateClient, error) {
	return nil, nil
}

func (msg MsgCreateClient) Route() string {
	return host.RouterKey
}

func (msg MsgCreateClient) MsgType() string {
	return TypeMsgCreateClient
}

func (msg MsgCreateClient) ValidateBasic() abcitypes.Error {
	if err := msg.ClientState.Validate(); err != nil {
		return ErrorClientState(DefaultCodespace, err)
	}
	if msg.ClientState.ClientType() == exported.Localhost {
		return ErrorClientState(DefaultCodespace, errors.New("localhost client can only be created on chain initialization"))
	}
	if msg.ClientState.ClientType() != msg.ConsensusState.ClientType() {
		return ErrorClientState(DefaultCodespace, errors.New("client type for client state and consensus state do not match"))
	}
	if err := ValidateClientType(msg.ClientState.ClientType()); err != nil {

	}
	if err := msg.ConsensusState.ValidateBasic(); err != nil {
		return ErrorConsensusState(DefaultCodespace, err)
	}
	return nil
}

func (msg MsgCreateClient) GetFromAddress() abcitypes.AccAddress {
	return abcitypes.HexToAddress(msg.Signer)
}

func (msg MsgCreateClient) Bytes() []byte {
	panic("IBC messages do not support amino")
}

// MsgCreateClientResponse defines the Msg/CreateClient response type.
type MsgCreateClientResponse struct {
}
