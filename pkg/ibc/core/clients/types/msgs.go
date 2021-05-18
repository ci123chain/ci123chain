package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
)

const (
	TypeMsgCreateClient		string = "create_client"
	TypeMsgUpdateClient       string = "update_client"

)

var _ sdk.Msg = &MsgCreateClient{}

type MsgCreateClient struct {
	// client_state
	ClientState exported.ClientState `json:"client_state,omitempty" yaml:"client_state"`
	// consensus_state
	ConsensusState exported.ConsensusState `json:"consensus_state,omitempty" yaml:"consensus_state"`
	// singer address
	Signer               string   `json:"signer,omitempty"`
}

func NewMsgCreateClient(
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
	signer sdk.AccAddress,
	) *MsgCreateClient {

	return &MsgCreateClient{
		ClientState: clientState,
		ConsensusState: consensusState,
		Signer: signer.String(),
	}
}

func (msg MsgCreateClient) Route() string {
	return host.RouterKey
}

func (msg MsgCreateClient) MsgType() string {
	return TypeMsgCreateClient
}

func (msg MsgCreateClient) ValidateBasic() error {
	if err := msg.ClientState.Validate(); err != nil {

		return ErrInvalidClient(err.Error())
	}
	if msg.ClientState.ClientType() == exported.Localhost {
		return ErrInvalidClient("localhost client can only be created on chain initialization")
	}
	if msg.ClientState.ClientType() != msg.ConsensusState.ClientType() {
		return ErrInvalidClient("client type for client state and consensus state do not match")
	}
	if err := ValidateClientType(msg.ClientState.ClientType()); err != nil {
		return ErrInvalidClient("client type does not meet naming constraints")
	}
	return msg.ConsensusState.ValidateBasic()
}

func (msg MsgCreateClient) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(msg.Signer)
}

func (msg MsgCreateClient) Bytes() []byte {
	panic("IBC messages do not support amino")
}

// MsgCreateClientResponse defines the Msg/CreateClient response types.
type MsgCreateClientResponse struct {
}

var _ sdk.Msg = &MsgUpdateClient{}
// MsgUpdateClient defines an sdk.Msg to update a IBC client state using
// the given header.
type MsgUpdateClient struct {
	// client unique identifier
	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty" yaml:"client_id"`
	// header to update the light client
	Header exported.Header `protobuf:"bytes,2,opt,name=header,proto3" json:"header,omitempty"`
	// signer address
	Signer string `protobuf:"bytes,3,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m MsgUpdateClient) Route() string {
	return host.RouterKey
}

func (m MsgUpdateClient) MsgType() string {
	return TypeMsgUpdateClient
}

func (msg MsgUpdateClient) ValidateBasic() error {

	if err := msg.Header.ValidateBasic(); err != nil {
		return ErrInvalidParam(err.Error())
	}
	if msg.ClientId == exported.Localhost {
		return ErrInvalidClient("localhost client is only updated on ABCI BeginBlock")
	}
	return host.ClientIdentifierValidator(msg.ClientId)
}

func (m MsgUpdateClient) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgUpdateClient) Bytes() []byte {
	panic("implement me")
}

// NewMsgUpdateClient creates a new MsgUpdateClient instance
//nolint:interfacer
func NewMsgUpdateClient(id string, header exported.Header, signer sdk.AccAddress) *MsgUpdateClient {
	return &MsgUpdateClient{
		ClientId: id,
		Header:   header,
		Signer:   signer.String(),
	}
}


// MsgUpdateClientResponse defines the Msg/UpdateClient response type.
type MsgUpdateClientResponse struct {
}
