package transfer

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)



type MsgTransfer struct {
	FromAddress sdk.AccAddress  `json:"from"`
	To     		sdk.AccAddress  `json:"to"`
	Amount 		sdk.Coins        `json:"amount"`
	FabricMode 	bool         	`json:"fabric_mode"`
}

func NewMsgTransfer(from, to sdk.AccAddress, amount sdk.Coins, isFabric bool ) *MsgTransfer {
	msg := &MsgTransfer{
		FromAddress: 	from,
		To: 			to,
		Amount: 		amount,
		FabricMode: 	isFabric,
	}
	return msg
}

func (msg *MsgTransfer) ValidateBasic() error {
	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty to address")
	}
	return nil
}

func (msg *MsgTransfer) Route() string { return RouteKey }

func (msg *MsgTransfer) MsgType() string { return "transfer"}

func (msg *MsgTransfer) Bytes() []byte {
	bytes, err := transferCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgTransfer) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}