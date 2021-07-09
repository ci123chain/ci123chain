package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/go-amino"
)


var IBCClientCodec *codec.Codec

func init() {
	IBCClientCodec = codec.New()
	RegisterCodec(IBCClientCodec)
	codec.RegisterCrypto(IBCClientCodec)

	IBCClientCodec.Seal()
}
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ClientState)(nil), &amino.InterfaceOptions{AlwaysDisambiguate: true})
	cdc.RegisterInterface((*exported.ConsensusState)(nil), nil)
	cdc.RegisterInterface((*exported.Header)(nil), &amino.InterfaceOptions{AlwaysDisambiguate: true})

	cdc.RegisterConcrete(&MsgCreateClient{}, "ibcclient/MsgCreateClient", nil)
	cdc.RegisterConcrete(&MsgUpdateClient{}, "ibcclient/MsgUpdateClient", nil)
	cdc.RegisterConcrete(&Height{}, "ibcclient/Height", nil)

	cdc.RegisterConcrete(&GenesisState{}, "ibcclient/GenesisState", nil)
	cdc.RegisterConcrete(&ConsensusStateWithHeight{}, "ibcclient/ConsensusStateWithHeight", nil)
	cdc.RegisterConcrete(&ClientConsensusStates{}, "ibcclient/ClientConsensusStates", nil)
	cdc.RegisterConcrete(&IdentifiedClientState{}, "ibcclient/IdentifiedClientState", nil)

	cdc.RegisterConcrete(&QueryClientStateResponse{}, "ibcclient/QueryClientStateResponse", nil)
	cdc.RegisterConcrete(&QueryClientStatesResponse{}, "ibcclient/QueryClientStatesResponse", nil)

}


func SetBinary(registry codectypes.InterfaceRegistry) {
	SubModuleCdc = codec.NewProtoCodec(registry)
}

var (
	// SubModuleCdc references the global x/ibc/core/03-connection module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to x/ibc/core/03-connection and
	// defined at the application level.
	//var Marshaler codec.Marshaler
	SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

// RegisterInterfaces registers the client interfaces to protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"ibc.core.client.v1.ClientState",
		(*exported.ClientState)(nil),
	)
	registry.RegisterInterface(
		"ibc.core.client.v1.ConsensusState",
		(*exported.ConsensusState)(nil),
	)
	registry.RegisterInterface(
		"ibc.core.client.v1.Header",
		(*exported.Header)(nil),
	)
	registry.RegisterInterface(
		"ibc.core.client.v1.Height",
		(*exported.Height)(nil),
		&Height{},
	)
	////registry.RegisterInterface(
	////	"ibc.core.client.v1.Misbehaviour",
	////	(*exported.Misbehaviour)(nil),
	////)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateClient{},
		&MsgUpdateClient{},
		//&MsgUpgradeClient{},
		//&MsgSubmitMisbehaviour{},
	)

	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// PackClientState constructs a new Any packed with the given client state value. It returns
// an error if the client state can't be casted to a protobuf message or if the concrete
// implemention is not registered to the protobuf codec.
func PackClientState(clientState exported.ClientState) (*codectypes.Any, error) {
	msg, ok := clientState.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", clientState)
	}

	anyClientState, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyClientState, nil
}

// UnpackClientState unpacks an Any into a ClientState. It returns an error if the
// client state can't be unpacked into a ClientState.
func UnpackClientState(any *codectypes.Any) (exported.ClientState, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	clientState, ok := any.GetCachedValue().(exported.ClientState)
	if ok {
		return clientState, nil
	}
	err := SubModuleCdc.UnpackAny(any, &clientState)

	return clientState, err
}

// PackConsensusState constructs a new Any packed with the given consensus state value. It returns
// an error if the consensus state can't be casted to a protobuf message or if the concrete
// implemention is not registered to the protobuf codec.
func PackConsensusState(consensusState exported.ConsensusState) (*codectypes.Any, error) {
	msg, ok := consensusState.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", consensusState)
	}

	anyConsensusState, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyConsensusState, nil
}

// MustPackConsensusState calls PackConsensusState and panics on error.
func MustPackConsensusState(consensusState exported.ConsensusState) *codectypes.Any {
	anyConsensusState, err := PackConsensusState(consensusState)
	if err != nil {
		panic(err)
	}

	return anyConsensusState
}

// UnpackConsensusState unpacks an Any into a ConsensusState. It returns an error if the
// consensus state can't be unpacked into a ConsensusState.
func UnpackConsensusState(any *codectypes.Any) (exported.ConsensusState, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	consensusState, ok := any.GetCachedValue().(exported.ConsensusState)
	if ok {
		return consensusState, nil
	}

	err := codectypes.AminoUnpacker{Cdc: IBCClientCodec}.UnpackAny(any, &consensusState)

	return consensusState, err
}

// PackHeader constructs a new Any packed with the given header value. It returns
// an error if the header can't be casted to a protobuf message or if the concrete
// implemention is not registered to the protobuf codec.
func PackHeader(header exported.Header) (*codectypes.Any, error) {
	msg, ok := header.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", header)
	}

	anyHeader, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyHeader, nil
}

// UnpackHeader unpacks an Any into a Header. It returns an error if the
// consensus state can't be unpacked into a Header.
func UnpackHeader(any *codectypes.Any) (exported.Header, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	header, ok := any.GetCachedValue().(exported.Header)
	if ok {
		return header, nil
	}

	err := codectypes.AminoUnpacker{Cdc: IBCClientCodec}.UnpackAny(any, &header)

	return header, err
}
