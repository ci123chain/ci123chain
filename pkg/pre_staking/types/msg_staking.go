package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"math/big"
)

type MsgStaking struct {
	FromAddress    sdk.AccAddress   `json:"from_address"`
	Delegator      sdk.AccAddress   `json:"delegator"`
	Validator      sdk.AccAddress   `json:"validator"`
	VaultID        string         `json:"vault_id"`
}

func NewMsgStaking(from sdk.AccAddress, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress,
	id string) *MsgStaking {
		return &MsgStaking{
			FromAddress: from,
			Delegator:   delegatorAddr,
			Validator:   validatorAddr,
			VaultID:      id,
		}
}

func (msg *MsgStaking) ValidateBasic() error {
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
	id, ok := new(big.Int).SetString(msg.VaultID, 10)
	if !ok {
		return ErrInvalidVaultID
	}
	if id.Uint64() < 1 {
		return ErrInvalidVaultID
	}

	return nil
}

func (msg *MsgStaking) Bytes() []byte {
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgStaking) Route() string {return ModuleName}
func (msg *MsgStaking) MsgType() string {return "delegate"}
func (msg *MsgStaking) GetFromAddress() sdk.AccAddress { return msg.FromAddress}