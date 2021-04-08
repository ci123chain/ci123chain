package types

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	sdkerrors "github.com/pkg/errors"
)

// Channel defines pipeline for exactly-once packet delivery between specific
// modules on separate blockchains, which has at least one end capable of
// sending packets and one end capable of receiving packets.
type Channel struct {
	// current state of the channel end
	State State `protobuf:"varint,1,opt,name=state,proto3,enum=ibc.collactor.channel.v1.State" json:"state,omitempty"`
	// whether the channel is ordered or unordered
	Ordering Order `protobuf:"varint,2,opt,name=ordering,proto3,enum=ibc.collactor.channel.v1.Order" json:"ordering,omitempty"`
	// counterparty channel end
	Counterparty Counterparty `protobuf:"bytes,3,opt,name=counterparty,proto3" json:"counterparty"`
	// list of connection identifiers, in order, along which packets sent on
	// this channel will travel
	ConnectionHops []string `protobuf:"bytes,4,rep,name=connection_hops,json=connectionHops,proto3" json:"connection_hops,omitempty" yaml:"connection_hops"`
	// opaque channel version, which is agreed upon during the handshake
	Version string `protobuf:"bytes,5,opt,name=version,proto3" json:"version,omitempty"`
}



// Order defines if a channel is ORDERED or UNORDERED
type Order int32

const (
	// zero-value for channel ordering
	NONE Order = 0
	// packets can be delivered in any order, which may differ from the order in
	// which they were sent.
	UNORDERED Order = 1
	// packets are delivered exactly in the order which they were sent
	ORDERED Order = 2
)

var Order_name = map[int32]string{
	0: "ORDER_NONE_UNSPECIFIED",
	1: "ORDER_UNORDERED",
	2: "ORDER_ORDERED",
}

var Order_value = map[string]int32{
	"ORDER_NONE_UNSPECIFIED": 0,
	"ORDER_UNORDERED":        1,
	"ORDER_ORDERED":          2,
}

func (x Order) String() string {
	return Order_name[int32(x)]
}

type State int32

const (
	// Default State
	UNINITIALIZED State = 0
	// A channel has just started the opening handshake.
	INIT State = 1
	// A channel has acknowledged the handshake step on the counterparty chain.
	TRYOPEN State = 2
	// A channel has completed the handshake. Open channels are
	// ready to send and receive packets.
	OPEN State = 3
	// A channel has been closed and can no longer be used to send or receive
	// packets.
	CLOSED State = 4
)

var State_name = map[int32]string{
	0: "STATE_UNINITIALIZED_UNSPECIFIED",
	1: "STATE_INIT",
	2: "STATE_TRYOPEN",
	3: "STATE_OPEN",
	4: "STATE_CLOSED",
}

var State_value = map[string]int32{
	"STATE_UNINITIALIZED_UNSPECIFIED": 0,
	"STATE_INIT":                      1,
	"STATE_TRYOPEN":                   2,
	"STATE_OPEN":                      3,
	"STATE_CLOSED":                    4,
}

func (x State) String() string {
	return Order_name[int32((x))]
}


// Counterparty defines a channel end counterparty
type Counterparty struct {
	// port on the counterparty chain which owns the other end of the channel.
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty" yaml:"port_id"`
	// channel end on the counterparty chain
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty" yaml:"channel_id"`
}





var (
	_ exported.ChannelI             = (*Channel)(nil)
	_ exported.CounterpartyChannelI = (*Counterparty)(nil)
)

// NewChannel creates a new Channel instance
func NewChannel(
	state State, ordering Order, counterparty Counterparty,
	hops []string, version string,
) Channel {
	return Channel{
		State:          state,
		Ordering:       ordering,
		Counterparty:   counterparty,
		ConnectionHops: hops,
		Version:        version,
	}
}

// GetState implements Channel interface.
func (ch Channel) GetState() int32 {
	return int32(ch.State)
}

// GetOrdering implements Channel interface.
func (ch Channel) GetOrdering() int32 {
	return int32(ch.Ordering)
}

// GetCounterparty implements Channel interface.
func (ch Channel) GetCounterparty() exported.CounterpartyChannelI {
	return ch.Counterparty
}

// GetConnectionHops implements Channel interface.
func (ch Channel) GetConnectionHops() []string {
	return ch.ConnectionHops
}

// GetVersion implements Channel interface.
func (ch Channel) GetVersion() string {
	return ch.Version
}

// ValidateBasic performs a basic validation of the channel fields
func (ch Channel) ValidateBasic() error {
	if ch.State == UNINITIALIZED {
		return ErrInvalidChannelState
	}
	if !(ch.Ordering == ORDERED || ch.Ordering == UNORDERED) {
		return sdkerrors.Wrap(ErrInvalidChannelOrdering, ch.Ordering.String())
	}
	if len(ch.ConnectionHops) != 1 {
		return sdkerrors.Wrap(
			ErrTooManyConnectionHops,
			"current IBC version only supports one connection hop",
		)
	}
	if err := host.ConnectionIdentifierValidator(ch.ConnectionHops[0]); err != nil {
		return sdkerrors.Wrap(err, "invalid connection hop ID")
	}
	return ch.Counterparty.ValidateBasic()
}

// NewCounterparty returns a new Counterparty instance
func NewCounterparty(portID, channelID string) Counterparty {
	return Counterparty{
		PortId:    portID,
		ChannelId: channelID,
	}
}

// GetPortID implements CounterpartyChannelI interface
func (c Counterparty) GetPortID() string {
	return c.PortId
}

// GetChannelID implements CounterpartyChannelI interface
func (c Counterparty) GetChannelID() string {
	return c.ChannelId
}

// ValidateBasic performs a basic validation check of the identifiers
func (c Counterparty) ValidateBasic() error {
	if err := host.PortIdentifierValidator(c.PortId); err != nil {
		return sdkerrors.Wrap(err, "invalid counterparty port ID")
	}
	if c.ChannelId != "" {
		if err := host.ChannelIdentifierValidator(c.ChannelId); err != nil {
			return sdkerrors.Wrap(err, "invalid counterparty channel ID")
		}
	}
	return nil
}

type isAcknowledgement_Response interface {
	isAcknowledgement_Response()
	MarshalTo([]byte) (int, error)
	Size() int
}


type Acknowledgement_Result struct {
	Result []byte `protobuf:"bytes,21,opt,name=result,proto3,oneof" json:"result,omitempty"`
}
type Acknowledgement_Error struct {
	Error string `protobuf:"bytes,22,opt,name=error,proto3,oneof" json:"error,omitempty"`
}
//
//func (*Acknowledgement_Result) isAcknowledgement_Response() {}
//func (*Acknowledgement_Error) isAcknowledgement_Response()  {}


type Acknowledgement struct {
	Response interface{}
}

func (ack Acknowledgement) String() string {
	res, _ := json.Marshal(ack)
	return string(res)
}
// GetBytes is a helper for serialising acknowledgements
func (ack Acknowledgement) GetBytes() []byte {
	return ChannelCdc.MustMarshalJSON(ack)
}

// NewResultAcknowledgement returns a new instance of Acknowledgement using an Acknowledgement_Result
// type in the Response field.
func NewResultAcknowledgement(result []byte) Acknowledgement {
	return Acknowledgement{
		Response: &Acknowledgement_Result{
			Result: result,
		},
	}
}

// NewErrorAcknowledgement returns a new instance of Acknowledgement using an Acknowledgement_Error
// type in the Response field.
func NewErrorAcknowledgement(err string) Acknowledgement {
	return Acknowledgement{
		Response: &Acknowledgement_Error{
			Error: err,
		},
	}
}
