package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"time"
)

type MsgStakingDirect struct {
	FromAddress    sdk.AccAddress   `json:"from_address"`
	Delegator      sdk.AccAddress   `json:"delegator"`
	Validator      sdk.AccAddress   `json:"validator"`
	Amount    	   sdk.Coin		    `json:"amount"`
	DelegateTime   time.Duration    `json:"delegate_time"`
}

func NewMsgStakingDirect(from sdk.AccAddress, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress,
	amount sdk.Coin, dt time.Duration) *MsgStakingDirect {
	return &MsgStakingDirect{
		FromAddress: from,
		Delegator:   delegatorAddr,
		Validator:   validatorAddr,
		Amount:      amount,
		DelegateTime: dt,
	}
}

func (msg *MsgStakingDirect) ValidateBasic() error {
	if msg.Delegator.Empty() {
		return ErrInvalidDelegatorAddress
	}
	if msg.Validator.Empty() {
		return ErrInvalidValidatorAddress
	}
	if !msg.FromAddress.Equal(msg.Delegator) {
		return ErrFromNotEqualDelegator
	}
	//if msg.VaultID.Cmp(new(big.Int).SetUint64(0)) <= 1 {
	//	return ErrInvalidVaultID
	//}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	}
	if msg.DelegateTime.Seconds() <= (time.Second * 60).Seconds(){
		return sdkerrors.Wrap(sdkerrors.ErrParams, "the time should longer than 168h(1 week)")
	}

	return nil
}

func (msg *MsgStakingDirect) Bytes() []byte {
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgStakingDirect) Route() string {return ModuleName}
func (msg *MsgStakingDirect) MsgType() string {return "delegate-direct"}
func (msg *MsgStakingDirect) GetFromAddress() sdk.AccAddress { return msg.FromAddress}