package types

import (
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/types"
	"github.com/pkg/errors"
)

var _ exported.ConnectionI = (*ConnectionEnd)(nil)

type ConnectionEnd struct {
	// client associated with this connection.
	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty" yaml:"client_id"`
	// IBC version which can be utilised to determine encodings or protocols for
	// channels or packets utilising this connection.
	Versions []*Version `protobuf:"bytes,2,rep,name=versions,proto3" json:"versions,omitempty"`
	// current state of the connection end.
	State State `protobuf:"varint,3,opt,name=state,proto3,enum=ibc.core.connection.v1.State" json:"state,omitempty"`
	// counterparty chain associated with this connection.
	Counterparty Counterparty `protobuf:"bytes,4,opt,name=counterparty,proto3" json:"counterparty"`
	// delay period that must pass before a consensus state can be used for packet-verification
	// NOTE: delay period logic is only implemented by some clients.
	DelayPeriod uint64 `protobuf:"varint,5,opt,name=delay_period,json=delayPeriod,proto3" json:"delay_period,omitempty" yaml:"delay_period"`

}

// NewConnectionEnd creates a new ConnectionEnd instance.
func NewConnectionEnd(state State, clientID string, counterparty Counterparty, versions []*Version, delayPeriod uint64) ConnectionEnd {
	return ConnectionEnd{
		ClientId:     clientID,
		Versions:     versions,
		State:        state,
		Counterparty: counterparty,
		DelayPeriod:  delayPeriod,
	}
}

// GetState implements the Connection interface
func (c ConnectionEnd) GetState() int32 {
	return int32(c.State)
}

// GetClientID implements the Connection interface
func (c ConnectionEnd) GetClientID() string {
	return c.ClientId
}

// GetCounterparty implements the Connection interface
func (c ConnectionEnd) GetCounterparty() exported.CounterpartyConnectionI {
	return c.Counterparty
}

// GetVersions implements the Connection interface
func (c ConnectionEnd) GetVersions() []exported.Version {
	return VersionsToExported(c.Versions)
}

// GetDelayPeriod implements the Connection interface
func (c ConnectionEnd) GetDelayPeriod() uint64 {
	return c.DelayPeriod
}

// ValidateBasic implements the Connection interface.
// NOTE: the protocol supports that the connection and client IDs match the
// counterparty's.
func (c ConnectionEnd) ValidateBasic() error {
	if err := host.ClientIdentifierValidator(c.ClientId); err != nil {
		return errors.Wrap(err, "invalid client ID")
	}
	if len(c.Versions) == 0 {
		return errors.New("empty connection versions")
	}
	for _, version := range c.Versions {
		if err := ValidateVersion(version); err != nil {
			return err
		}
	}
	return c.Counterparty.ValidateBasic()
}

var _ exported.CounterpartyConnectionI = (*Counterparty)(nil)


// Counterparty defines the counterparty chain associated with a connection end.
type Counterparty struct {
	// identifies the client on the counterparty chain associated with a given
	// connection.
	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty" yaml:"client_id"`
	// identifies the connection end on the counterparty chain associated with a
	// given connection.
	ConnectionId string `protobuf:"bytes,2,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty" yaml:"connection_id"`
	// commitment merkle prefix of the counterparty chain.
	Prefix commitmenttypes.MerklePrefix `protobuf:"bytes,3,opt,name=prefix,proto3" json:"prefix"`
}

func NewCounterparty(clientID, connectionID string, prefix commitmenttypes.MerklePrefix) Counterparty {
	return Counterparty{
		ClientId: clientID,
		ConnectionId: connectionID,
		Prefix: prefix,
	}
}


// GetClientID implements the CounterpartyConnectionI interface
func (c Counterparty) GetClientID() string {
	return c.ClientId
}

// GetConnectionID implements the CounterpartyConnectionI interface
func (c Counterparty) GetConnectionID() string {
	return c.ConnectionId
}

// GetPrefix implements the CounterpartyConnectionI interface
func (c Counterparty) GetPrefix() exported.Prefix {
	return &c.Prefix
}

// ValidateBasic performs a basic validation check of the identifiers and prefix
func (c Counterparty) ValidateBasic() error {
	if c.ConnectionId != "" {
		if err := host.ConnectionIdentifierValidator(c.ConnectionId); err != nil {
			return types.ErrorCounterpartyConnectionID(types.DefaultCodespace, err)
		}
	}
	if err := host.ClientIdentifierValidator(c.ClientId); err != nil {
		return types.ErrorCounterpartyConnectionID(types.DefaultCodespace, err)
	}
	if c.Prefix.Empty() {
		return types.ErrorCounterpartyPrefix(types.DefaultCodespace, errors.New("counterparty prefix cannot be empty"))
	}
	return nil
}

type State int32

const (
	// Default State
	UNINITIALIZED State = 0
	// A connection end has just started the opening handshake.
	INIT State = 1
	// A connection end has acknowledged the handshake step on the counterparty
	// chain.
	TRYOPEN State = 2
	// A connection end has completed the handshake.
	OPEN State = 3
)

var State_name = map[int32]string{
	0: "STATE_UNINITIALIZED_UNSPECIFIED",
	1: "STATE_INIT",
	2: "STATE_TRYOPEN",
	3: "STATE_OPEN",
}

var State_value = map[string]int32{
	"STATE_UNINITIALIZED_UNSPECIFIED": 0,
	"STATE_INIT":                      1,
	"STATE_TRYOPEN":                   2,
	"STATE_OPEN":                      3,
}


func (x State) String() string {
	return State_name[int32(x)]
}



