package types

import (
	"fmt"
	paramtypes "github.com/ci123chain/ci123chain/pkg/params/types"
)


const (
	// DefaultSendEnabled enabled
	DefaultSendEnabled = true
	// DefaultReceiveEnabled enabled
	DefaultReceiveEnabled = true
)
var (
	// KeySendEnabled is store's key for SendEnabled Params
	KeySendEnabled = []byte("SendEnabled")
	// KeyReceiveEnabled is store's key for ReceiveEnabled Params
	KeyReceiveEnabled = []byte("ReceiveEnabled")
)


type Params struct {
	// send_enabled enables or disables all cross-chain token transfers from this
	// chain.
	SendEnabled bool `protobuf:"varint,1,opt,name=send_enabled,json=sendEnabled,proto3" json:"send_enabled,omitempty" yaml:"send_enabled"`
	// receive_enabled enables or disables all cross-chain token transfers to this
	// chain.
	ReceiveEnabled bool `protobuf:"varint,2,opt,name=receive_enabled,json=receiveEnabled,proto3" json:"receive_enabled,omitempty" yaml:"receive_enabled"`
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the ibc transfer module
func NewParams(enableSend, enableReceive bool) Params {
	return Params{
		SendEnabled:    enableSend,
		ReceiveEnabled: enableReceive,
	}
}

// Validate all ibc-transfer module parameters
func (p Params) Validate() error {
	if err := validateEnabled(p.SendEnabled); err != nil {
		return err
	}

	return validateEnabled(p.ReceiveEnabled)
}

func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeySendEnabled, p.SendEnabled, validateEnabled),
		paramtypes.NewParamSetPair(KeyReceiveEnabled, p.ReceiveEnabled, validateEnabled),
	}
}

// DefaultParams is the default parameter configuration for the ibc-transfer module
func DefaultParams() Params {
	return NewParams(DefaultSendEnabled, DefaultReceiveEnabled)
}

