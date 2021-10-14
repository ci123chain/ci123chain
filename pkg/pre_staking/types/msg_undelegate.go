package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)


type MsgUndelegate struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	VaultID     string       `json:"vault_id"`
}

func NewMsgUndelegate(from sdk.AccAddress, id string) *MsgUndelegate {
	return &MsgUndelegate{
		FromAddress: from,
		VaultID:      id,
	}
}

func (msg *MsgUndelegate) Route() string { return ModuleName }

func (msg *MsgUndelegate) MsgType() string { return "undelegate" }

func (msg *MsgUndelegate) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "from address can not empty")
	}
	return nil
}


func (msg *MsgUndelegate) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

func (msg *MsgUndelegate) Bytes() []byte{
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}