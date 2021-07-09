package types

import (
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreateClient		string = "create_client"
	TypeMsgUpdateClient       string = "update_client"

)

var _ sdk.Msg = &MsgCreateClient{}
var _ cosmosSdk.Msg = &MsgCreateClient{}

//type MsgCreateClient struct {
//	// client_state
//	ClientState exported.ClientState `json:"client_state,omitempty" yaml:"client_state"`
//	// consensus_state
//	ConsensusState exported.ConsensusState `json:"consensus_state,omitempty" yaml:"consensus_state"`
//	// singer address
//	Signer               string   `json:"signer,omitempty"`
//}
//func (m *MsgCreateClient) Reset()         { *m = MsgCreateClient{} }
//func (m *MsgCreateClient) String() string { return proto.CompactTextString(m) }
//func (*MsgCreateClient) ProtoMessage()    {}
func NewMsgCreateClient(
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
	signer string,
	) (*MsgCreateClient, error) {

	anyClientState, err := PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	anyConsensusState, err := PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}

	return &MsgCreateClient{
		ClientState: anyClientState,
		ConsensusState: anyConsensusState,
		Signer: signer,
	}, nil
	//return &MsgCreateClient{
	//	ClientState: anyClientState,
	//	ConsensusState: anyConsensusState,
	//	Signer: signer.String(),
	//}, nil
}

func (msg MsgCreateClient) Route() string {
	return host.RouterKey
}

func (msg MsgCreateClient) MsgType() string {
	return TypeMsgCreateClient
}

func (msg MsgCreateClient) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	clientState, err := UnpackClientState(msg.ClientState)
	if err != nil {
		return err
	}
	if err := clientState.Validate(); err != nil {
		return err
	}
	if clientState.ClientType() == exported.Localhost {
		return sdkerrors.Wrap(ErrInvalidClient, "localhost client can only be created on chain initialization")
	}
	consensusState, err := UnpackConsensusState(msg.ConsensusState)
	if err != nil {
		return err
	}
	if clientState.ClientType() != consensusState.ClientType() {
		return sdkerrors.Wrap(ErrInvalidClientType, "client type for client state and consensus state do not match")
	}
	if err := ValidateClientType(clientState.ClientType()); err != nil {
		return sdkerrors.Wrap(err, "client type does not meet naming constraints")
	}
	return consensusState.ValidateBasic()
}

func (msg MsgCreateClient) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(msg.Signer)
}

func (msg MsgCreateClient) Bytes() []byte {
	panic("IBC messages do not support amino")
}

func (msg MsgCreateClient) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

func (msg MsgCreateClient) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(msg.Signer).Bytes()}
}

func (msg MsgCreateClient) Type() string {
	return TypeMsgCreateClient
}

// MsgCreateClientResponse defines the Msg/CreateClient response types.
//type MsgCreateClientResponse struct {
//}

var _ sdk.Msg = &MsgUpdateClient{}
var _ cosmosSdk.Msg = &MsgUpdateClient{}
// MsgUpdateClient defines an sdk.Msg to update a IBC client state using
// the given header.
//type MsgUpdateClient struct {
//	// client unique identifier
//	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty" yaml:"client_id"`
//	// header to update the light client
//	Header exported.Header `protobuf:"bytes,2,opt,name=header,proto3" json:"header,omitempty"`
//	// signer address
//	Signer string `protobuf:"bytes,3,opt,name=signer,proto3" json:"signer,omitempty"`
//}

func (m MsgUpdateClient) Route() string {
	return host.RouterKey
}

func (m MsgUpdateClient) MsgType() string {
	return TypeMsgUpdateClient
}

func (msg MsgUpdateClient) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	header, err := UnpackHeader(msg.Header)
	if err != nil {
		return err
	}
	if err := header.ValidateBasic(); err != nil {
		return err
	}
	if msg.ClientId == exported.Localhost {
		return sdkerrors.Wrap(ErrInvalidClient, "localhost client is only updated on ABCI BeginBlock")
	}
	return host.ClientIdentifierValidator(msg.ClientId)
}

func (m MsgUpdateClient) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgUpdateClient) Bytes() []byte {
	panic("implement me")
}

func (m MsgUpdateClient) Type() string {
	return TypeMsgUpdateClient
}

func (m MsgUpdateClient) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgUpdateClient) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(m.Signer).Bytes()}
}

// NewMsgUpdateClient creates a new MsgUpdateClient instance
//nolint:interfacer
func NewMsgUpdateClient(id string, header exported.Header, signer string) (*MsgUpdateClient, error) {
	anyHeader, err := PackHeader(header)
	if err != nil {
		return nil, err
	}

	return &MsgUpdateClient{
		ClientId: id,
		Header:   anyHeader,
		Signer:   signer,
	}, nil
}


// MsgUpdateClientResponse defines the Msg/UpdateClient response type.
//type MsgUpdateClientResponse struct {
//}

var _ codectypes.UnpackInterfacesMessage = MsgCreateClient{}
var _ codectypes.UnpackInterfacesMessage = MsgUpdateClient{}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateClient) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var clientState exported.ClientState
	err := unpacker.UnpackAny(msg.ClientState, &clientState)
	if err != nil {
		return err
	}

	var consensusState exported.ConsensusState
	return unpacker.UnpackAny(msg.ConsensusState, &consensusState)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgUpdateClient) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var header exported.Header
	return unpacker.UnpackAny(msg.Header, &header)
}