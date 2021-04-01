package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgUpgrade struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Type      	string   		`json:"types"`
	Height    	int64    		`json:"height"`
	Name      	string   		`json:"name"`
}

func NewMsgUpgrade(from sdk.AccAddress, t, name string, height int64) *MsgUpgrade{
	return &MsgUpgrade{
		FromAddress: from,
		Type:t,
		Height:height,
		Name:name,
	}
}

func (msg *MsgUpgrade) Route() string { return RouteKey }

func (msg *MsgUpgrade) MsgType() string { return "upgrade"}

func (msg *MsgUpgrade) ValidateBasic() sdk.Error{
	if len(msg.Type) == 0 {
		return ErrCheckParams(DefaultCodespace, "types is invalid")
	}
	if msg.Height < 0 {
		return ErrCheckParams(DefaultCodespace, "height is invalid")
	}
	if len(msg.Name) == 0 {
		return ErrCheckParams(DefaultCodespace, "name is invalid")
	}
	return nil
}

func (msg *MsgUpgrade) Bytes() []byte {
	bytes, err := ModuleCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgUpgrade) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}