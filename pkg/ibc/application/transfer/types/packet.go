package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/pkg/errors"
	"strings"
)


// NewFungibleTokenPacketData contructs a new FungibleTokenPacketData instance
func NewFungibleTokenPacketData(
	denom string, amount uint64,
	sender, receiver string,
) FungibleTokenPacketData {
	return FungibleTokenPacketData{
		Denom:    denom,
		Amount:   amount,
		Sender:   sender,
		Receiver: receiver,
	}
}


// ValidateBasic is used for validating the token transfer.
// NOTE: The addresses formats are not validated as the sender and recipient can have different
// formats defined by their corresponding chains that are not known to IBC.
func (ftpd FungibleTokenPacketData) ValidateBasic() error {
	if ftpd.Amount == 0 {
		return errors.Wrap(ErrInvalidAmount, "amount cannot be 0")
	}
	if strings.TrimSpace(ftpd.Sender) == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "sender address cannot be blank")
	}
	if strings.TrimSpace(ftpd.Receiver) == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "receiver address cannot be blank")
	}
	//return ValidatePrefixedDenom(ftpd.Denom)
	return nil
}

// GetBytes is a helper for serialising
func (ftpd FungibleTokenPacketData) GetBytes() []byte {
	bytes, err := IBCTransferCdc.MarshalJSON(ftpd)
	if err != nil {
		panic(err)
	}
	return bytes
}
