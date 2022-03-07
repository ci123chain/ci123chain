package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	proto "github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	_ sdk.Msg = &MsgSetOrchestratorAddress{}
	_ sdk.Msg = &MsgValsetConfirm{}
	_ sdk.Msg = &MsgSendToEth{}
	_ sdk.Msg = &MsgRequestBatch{}
	_ sdk.Msg = &MsgSend721ToEth{}
	_ sdk.Msg = &MsgRequest721Batch{}
	_ sdk.Msg = &MsgConfirmBatch{}
	_ sdk.Msg = &MsgConfirm721Batch{}
	_ sdk.Msg = &MsgERC20DeployedClaim{}
	_ sdk.Msg = &MsgERC721DeployedClaim{}
	_ sdk.Msg = &MsgConfirmLogicCall{}
	_ sdk.Msg = &MsgLogicCallExecutedClaim{}
	_ sdk.Msg = &MsgDepositClaim{}
	_ sdk.Msg = &MsgDeposit721Claim{}
	_ sdk.Msg = &MsgWithdrawClaim{}
	_ sdk.Msg = &MsgWithdraw721Claim{}
)

type MsgSend721ToEth struct {
	Sender    string     `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	EthDest   string     `protobuf:"bytes,2,opt,name=eth_dest,json=ethDest,proto3" json:"eth_dest,omitempty"`
	Amount    sdk.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
	BridgeFee sdk.Coin `protobuf:"bytes,4,opt,name=bridge_fee,json=bridgeFee,proto3" json:"bridge_fee"`
}

type MsgSend721ToEthResponse struct {
}

func (m *MsgSend721ToEthResponse) Reset()         { *m = MsgSend721ToEthResponse{} }
func (m *MsgSend721ToEthResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSend721ToEthResponse) ProtoMessage()    {}

type MsgRequest721Batch struct {
	Sender string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Denom  string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
}
func (m *MsgRequest721Batch) Reset()         { *m = MsgRequest721Batch{} }
func (m *MsgRequest721Batch) String() string { return proto.CompactTextString(m) }
func (*MsgRequest721Batch) ProtoMessage()    {}

type MsgRequest721BatchResponse struct {
}

func (m *MsgRequest721BatchResponse) Reset()         { *m = MsgRequest721BatchResponse{} }
func (m *MsgRequest721BatchResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRequest721BatchResponse) ProtoMessage()    {}

type MsgConfirm721Batch struct {
	Nonce         uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	TokenContract string `protobuf:"bytes,2,opt,name=token_contract,json=tokenContract,proto3" json:"token_contract,omitempty"`
	EthSigner     string `protobuf:"bytes,3,opt,name=eth_signer,json=ethSigner,proto3" json:"eth_signer,omitempty"`
	Orchestrator  string `protobuf:"bytes,4,opt,name=orchestrator,proto3" json:"orchestrator,omitempty"`
	Signature     string `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
}
func (m *MsgConfirm721Batch) Reset()         { *m = MsgConfirm721Batch{} }
func (m *MsgConfirm721Batch) String() string { return proto.CompactTextString(m) }
func (*MsgConfirm721Batch) ProtoMessage()    {}
type MsgConfirm721BatchResponse struct {
}

func (m *MsgConfirm721BatchResponse) Reset()         { *m = MsgConfirm721BatchResponse{} }
func (m *MsgConfirm721BatchResponse) String() string { return proto.CompactTextString(m) }
func (*MsgConfirm721BatchResponse) ProtoMessage()    {}

type MsgDeposit721Claim struct {
	EventNonce     uint64  `protobuf:"varint,1,opt,name=event_nonce,json=eventNonce,proto3" json:"event_nonce,omitempty"`
	BlockHeight    uint64  `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	TokenContract  string  `protobuf:"bytes,3,opt,name=token_contract,json=tokenContract,proto3" json:"token_contract,omitempty"`
	TokenName      string  `protobuf:"bytes,4,opt,name=token_name,json=tokenName,proto3" json:"token_name,omitempty"`
	TokenSymbol    string  `protobuf:"bytes,5,opt,name=token_symbol,json=tokenSymbol,proto3" json:"token_symbol,omitempty"`
	TokenID        sdk.Int `protobuf:"bytes,6,opt,name=token_id,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token_id"`
	EthereumSender string  `protobuf:"bytes,7,opt,name=ethereum_sender,json=ethereumSender,proto3" json:"ethereum_sender,omitempty"`
	CosmosReceiver string  `protobuf:"bytes,8,opt,name=cosmos_receiver,json=cosmosReceiver,proto3" json:"cosmos_receiver,omitempty"`
	Orchestrator   string  `protobuf:"bytes,9,opt,name=orchestrator,proto3" json:"orchestrator,omitempty"`
}

type MsgDeposit721ClaimResponse struct {
}

func (m *MsgDeposit721ClaimResponse) Reset()         { *m = MsgDeposit721ClaimResponse{} }
func (m *MsgDeposit721ClaimResponse) String() string { return proto.CompactTextString(m) }
func (*MsgDeposit721ClaimResponse) ProtoMessage()    {}

type MsgWithdraw721Claim struct {
	EventNonce    uint64 `protobuf:"varint,1,opt,name=event_nonce,json=eventNonce,proto3" json:"event_nonce,omitempty"`
	BlockHeight   uint64 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	BatchNonce    uint64 `protobuf:"varint,3,opt,name=batch_nonce,json=batchNonce,proto3" json:"batch_nonce,omitempty"`
	TokenContract string `protobuf:"bytes,4,opt,name=token_contract,json=tokenContract,proto3" json:"token_contract,omitempty"`
	Orchestrator  string `protobuf:"bytes,5,opt,name=orchestrator,proto3" json:"orchestrator,omitempty"`
}

type MsgERC721DeployedClaim struct {
	EventNonce    uint64 `protobuf:"varint,1,opt,name=event_nonce,json=eventNonce,proto3" json:"event_nonce,omitempty"`
	BlockHeight   uint64 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	CosmosDenom   string `protobuf:"bytes,3,opt,name=cosmos_denom,json=cosmosDenom,proto3" json:"cosmos_denom,omitempty"`
	TokenContract string `protobuf:"bytes,4,opt,name=token_contract,json=tokenContract,proto3" json:"token_contract,omitempty"`
	Name          string `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Symbol        string `protobuf:"bytes,6,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Orchestrator  string `protobuf:"bytes,7,opt,name=orchestrator,proto3" json:"orchestrator,omitempty"`
}

type MsgERC721DeployedClaimResponse struct {
}

func (m *MsgERC721DeployedClaimResponse) Reset()         { *m = MsgERC721DeployedClaimResponse{} }
func (m *MsgERC721DeployedClaimResponse) String() string { return proto.CompactTextString(m) }
func (*MsgERC721DeployedClaimResponse) ProtoMessage()    {}

type MsgValsetConfirmNonceClaim struct {
	ValsetNonce     uint64  `protobuf:"varint,1,opt,name=valset_nonce,json=valset_nonce,proto3" json:"valset_nonce,omitempty"`
	EventNonce		uint64  `protobuf:"varint,2,opt,name=event_nonce,json=event_nonce,proto3" json:"event_nonce,omitempty"`
	Orchestrator  string `protobuf:"bytes,3,opt,name=orchestrator,proto3" json:"orchestrator,omitempty"`
}

func (m *MsgValsetConfirmNonceClaim) Reset()         { *m = MsgValsetConfirmNonceClaim{} }
func (m *MsgValsetConfirmNonceClaim) String() string { return proto.CompactTextString(m) }
func (*MsgValsetConfirmNonceClaim) ProtoMessage()    {}

type MsgValsetConfirmNonceClaimResponse struct {
}

func (m *MsgValsetConfirmNonceClaimResponse) Reset()         { *m = MsgValsetConfirmNonceClaimResponse{} }
func (m *MsgValsetConfirmNonceClaimResponse) String() string { return proto.CompactTextString(m) }
func (*MsgValsetConfirmNonceClaimResponse) ProtoMessage()    {}

func (m *MsgValsetConfirmNonceClaim) GetEventNonce() uint64 {
	if m != nil {
		return m.EventNonce
	}
	return 0
}

func (m *MsgValsetConfirmNonceClaim) GetBlockHeight() uint64 {
	if m != nil {
		return 0
	}
	return 0
}

// GetType returns the type of the claim
func (e *MsgValsetConfirmNonceClaim) GetType() ClaimType {
	return CLAIM_TYPE_VALSET_CONFIRM_NONCE
}

// ValidateBasic performs stateless checks
func (e *MsgValsetConfirmNonceClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgValsetConfirmNonceClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgValsetConfirmNonceClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgValsetConfirmNonceClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgValsetConfirmNonceClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// MsgType should return the action
func (msg MsgValsetConfirmNonceClaim) MsgType() string { return "valset_confirm_nonce_claim" }

func (msg *MsgValsetConfirmNonceClaim) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// Route should return the name of the module
func (msg MsgValsetConfirmNonceClaim) Route() string { return RouterKey }

// ClaimHash implements BridgeDeposit.Hash
func (msg *MsgValsetConfirmNonceClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d", msg.ValsetNonce)
	return tmhash.Sum([]byte(path))
}

func (msg MsgValsetConfirmNonceClaim) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgSetOrchestratorAddress(val sdk.AccAddress, oper sdk.AccAddress, eth string) *MsgSetOrchestratorAddress {
	return &MsgSetOrchestratorAddress{
		Validator:    val.String(),
		Orchestrator: oper.String(),
		EthAddress:   eth,
	}
}

// Route should return the name of the module
func (msg *MsgSetOrchestratorAddress) Route() string { return RouterKey }

// MsgType should return the action
func (msg *MsgSetOrchestratorAddress) MsgType() string { return "set_operator_address" }

// ValidateBasic performs stateless checks
func (msg *MsgSetOrchestratorAddress) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}
	if _, err = sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgSetOrchestratorAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgSetOrchestratorAddress) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// GetFromAddress defines whose signature is required
func (msg *MsgSetOrchestratorAddress) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return acc
}

func (msg *MsgSetOrchestratorAddress) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// NewMsgValsetConfirm returns a new msgValsetConfirm
func NewMsgValsetConfirm(nonce uint64, ethAddress string, validator sdk.AccAddress, signature string) *MsgValsetConfirm {
	return &MsgValsetConfirm{
		Nonce:        nonce,
		Orchestrator: validator.String(),
		EthAddress:   ethAddress,
		Signature:    signature,
	}
}

// Route should return the name of the module
func (msg *MsgValsetConfirm) Route() string { return RouterKey }

// MsgType should return the action
func (msg *MsgValsetConfirm) MsgType() string { return "valset_confirm" }

// ValidateBasic performs stateless checks
func (msg *MsgValsetConfirm) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgValsetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgValsetConfirm) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg *MsgValsetConfirm) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return acc
}

func (msg *MsgValsetConfirm) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// NewMsgSendToEth returns a new msgSendToEth
func NewMsgSendToEth(sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEth {
	return &MsgSendToEth{
		Sender:    sender.String(),
		EthDest:   destAddress,
		Amount:    send,
		BridgeFee: bridgeFee,
	}
}

// Route should return the name of the module
func (msg MsgSendToEth) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgSendToEth) MsgType() string { return "send_to_eth" }

func (msg MsgSendToEth) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgSendToEth) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	//// fee and send must be of the same denom
	//if msg.Amount.Denom != msg.BridgeFee.Denom {
	//	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("fee and amount must be the same type %s != %s", msg.Amount.Denom, msg.BridgeFee.Denom))
	//}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if !msg.BridgeFee.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee")
	}
	if err := ValidateEthAddress(msg.EthDest); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO validate fee is sufficient, fixed fee to start
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSendToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSendToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

func (msg MsgSendToEth) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return acc
}

// Route should return the name of the module
func (msg MsgSend721ToEth) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgSend721ToEth) MsgType() string { return "send_721_to_eth" }

func (msg *MsgSend721ToEth) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgSend721ToEth) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	//// fee and send must be of the same denom
	//if msg.Amount.Denom != msg.BridgeFee.Denom {
	//	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("fee and amount must be the same type %s != %s", msg.Amount.Denom, msg.BridgeFee.Denom))
	//}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if !msg.BridgeFee.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee")
	}
	if err := ValidateEthAddress(msg.EthDest); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	// TODO validate fee is sufficient, fixed fee to start
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSend721ToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSend721ToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

func (msg MsgSend721ToEth) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return acc
}

// NewMsgRequestBatch returns a new msgRequestBatch
func NewMsgRequestBatch(orchestrator sdk.AccAddress) *MsgRequestBatch {
	return &MsgRequestBatch{
		Sender: orchestrator.String(),
	}
}

// Route should return the name of the module
func (msg MsgRequestBatch) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgRequestBatch) MsgType() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	return nil
}

func (msg *MsgRequestBatch) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

func (msg MsgRequestBatch) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return acc
}

// Route should return the name of the module
func (msg MsgRequest721Batch) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgRequest721Batch) MsgType() string { return "request_721_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequest721Batch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	return nil
}

func (msg *MsgRequest721Batch) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSignBytes encodes the message for signing
func (msg MsgRequest721Batch) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRequest721Batch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

func (msg MsgRequest721Batch) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return acc
}

// Route should return the name of the module
func (msg MsgConfirmBatch) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgConfirmBatch) MsgType() string { return "confirm_batch" }

// ValidateBasic performs stateless checks
func (msg MsgConfirmBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	if err := ValidateEthAddress(msg.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "token contract")
	}
	_, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	return nil
}

func (msg *MsgConfirmBatch) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSignBytes encodes the message for signing
func (msg MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg MsgConfirmBatch) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// Route should return the name of the module
func (msg MsgConfirm721Batch) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgConfirm721Batch) MsgType() string { return "confirm_721_batch" }

// ValidateBasic performs stateless checks
func (msg MsgConfirm721Batch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	if err := ValidateEthAddress(msg.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "token contract")
	}
	_, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	return nil
}

func (msg *MsgConfirm721Batch) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSignBytes encodes the message for signing
func (msg MsgConfirm721Batch) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConfirm721Batch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg MsgConfirm721Batch) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// Route should return the name of the module
func (msg MsgConfirmLogicCall) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgConfirmLogicCall) MsgType() string { return "confirm_logic" }

// ValidateBasic performs stateless checks
func (msg MsgConfirmLogicCall) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Orchestrator)
	}
	if err := ValidateEthAddress(msg.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	_, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.Signature)
	}
	_, err = hex.DecodeString(msg.InvalidationId)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", msg.InvalidationId)
	}
	return nil
}

func (msg *MsgConfirmLogicCall) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSignBytes encodes the message for signing
func (msg MsgConfirmLogicCall) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConfirmLogicCall) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg MsgConfirmLogicCall) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// EthereumClaim represents a claim on ethereum state
type EthereumClaim interface {
	// GetEventNonce All Ethereum claims that we relay from the Gravity contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// GetBlockHeight The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetBlockHeight() uint64
	// GetClaimer the delegate address of the claimer, for MsgDepositClaim and MsgWithdrawClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// GetType Which type of claim this is
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ EthereumClaim = &MsgDepositClaim{}
	_ EthereumClaim = &MsgDeposit721Claim{}
	_ EthereumClaim = &MsgWithdrawClaim{}
	_ EthereumClaim = &MsgWithdraw721Claim{}
	_ EthereumClaim = &MsgERC20DeployedClaim{}
	_ EthereumClaim = &MsgERC721DeployedClaim{}
	_ EthereumClaim = &MsgLogicCallExecutedClaim{}
)

// GetType returns the type of the claim
func (e *MsgDepositClaim) GetType() ClaimType {
	return CLAIM_TYPE_DEPOSIT
}

// ValidateBasic performs stateless checks
func (e *MsgDepositClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.CosmosReceiver)
	}
	if err := ValidateEthAddress(e.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDepositClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgDepositClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgDepositClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgDepositClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// MsgType should return the action
func (msg MsgDepositClaim) MsgType() string { return "deposit_claim" }

// Route should return the name of the module
func (msg MsgDepositClaim) Route() string { return RouterKey }

func (msg *MsgDepositClaim) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgDepositClaim) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// ClaimHash implements BridgeDeposit.Hash
func (msg *MsgDepositClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", msg.TokenContract, string(msg.EthereumSender), msg.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}

func (m *MsgDeposit721Claim) Reset()         { *m = MsgDeposit721Claim{} }
func (m *MsgDeposit721Claim) String() string { return proto.CompactTextString(m) }
func (*MsgDeposit721Claim) ProtoMessage()    {}

// GetType returns the type of the claim
func (e *MsgDeposit721Claim) GetType() ClaimType {
	return CLAIM_TYPE_DEPOSIT721
}

// ValidateBasic performs stateless checks
func (e *MsgDeposit721Claim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.CosmosReceiver)
	}
	if err := ValidateEthAddress(e.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeposit721Claim) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgDeposit721Claim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgDepositClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgDeposit721Claim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// MsgType should return the action
func (msg MsgDeposit721Claim) MsgType() string { return "deposit721_claim" }

// Route should return the name of the module
func (msg MsgDeposit721Claim) Route() string { return RouterKey }

func (msg *MsgDeposit721Claim) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgDeposit721Claim) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// ClaimHash implements BridgeDeposit.Hash
func (msg *MsgDeposit721Claim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", msg.TokenContract, string(msg.EthereumSender), msg.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}

func (m *MsgDeposit721Claim) GetEventNonce() uint64 {
	if m != nil {
		return m.EventNonce
	}
	return 0
}

func (m *MsgDeposit721Claim) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

const (
	TypeMsgWithdrawClaim = "withdraw_claim"
)

// GetType returns the claim type
func (e *MsgWithdrawClaim) GetType() ClaimType {
	return CLAIM_TYPE_WITHDRAW
}

// ValidateBasic performs stateless checks
func (e *MsgWithdrawClaim) ValidateBasic() error {
	if e.EventNonce == 0 {
		return fmt.Errorf("event_nonce == 0")
	}
	if e.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	return nil
}

// ClaimHash implements WithdrawBatch.Hash
func (b *MsgWithdrawClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", b.TokenContract, b.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (msg MsgWithdrawClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgWithdrawClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgWithdrawClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

func (msg MsgWithdrawClaim) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// GetSigners defines whose signature is required
func (msg MsgWithdrawClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgWithdrawClaim) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgWithdrawClaim) MsgType() string { return "withdraw_claim" }

func (msg *MsgWithdrawClaim) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

const (
	TypeMsgWithdraw721Claim = "withdraw_721_claim"
)

// GetType returns the claim type
func (e *MsgWithdraw721Claim) GetType() ClaimType {
	return CLAIM_TYPE_WITHDRAW
}

// ValidateBasic performs stateless checks
func (e *MsgWithdraw721Claim) ValidateBasic() error {
	if e.EventNonce == 0 {
		return fmt.Errorf("event_nonce == 0")
	}
	if e.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	return nil
}

// ClaimHash implements WithdrawBatch.Hash
func (b *MsgWithdraw721Claim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", b.TokenContract, b.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (msg MsgWithdraw721Claim) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgWithdraw721Claim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgWithdrawClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

func (msg MsgWithdraw721Claim) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// GetSigners defines whose signature is required
func (msg MsgWithdraw721Claim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (msg MsgWithdraw721Claim) Route() string { return RouterKey }

// MsgType should return the action
func (msg MsgWithdraw721Claim) MsgType() string { return "withdraw_721_claim" }

func (msg *MsgWithdraw721Claim) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

const (
	TypeMsgDepositClaim = "deposit_claim"
)

// GetType returns the type of the claim
func (e *MsgERC20DeployedClaim) GetType() ClaimType {
	return CLAIM_TYPE_ERC20_DEPLOYED
}

// ValidateBasic performs stateless checks
func (e *MsgERC20DeployedClaim) ValidateBasic() error {
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgERC20DeployedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgERC20DeployedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgERC20DeployedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// MsgType should return the action
func (msg MsgERC20DeployedClaim) MsgType() string { return "ERC20_deployed_claim" }

func (msg *MsgERC20DeployedClaim) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// Route should return the name of the module
func (msg MsgERC20DeployedClaim) Route() string { return RouterKey }

// ClaimHash implements BridgeDeposit.Hash
func (msg *MsgERC20DeployedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%s/%d/", msg.CosmosDenom, msg.TokenContract, msg.Name, msg.Symbol, msg.Decimals)
	return tmhash.Sum([]byte(path))
}

func (msg MsgERC20DeployedClaim) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

func (m *MsgERC721DeployedClaim) Reset()         { *m = MsgERC721DeployedClaim{} }
func (m *MsgERC721DeployedClaim) String() string { return proto.CompactTextString(m) }
func (*MsgERC721DeployedClaim) ProtoMessage()    {}

// GetType returns the type of the claim
func (e *MsgERC721DeployedClaim) GetType() ClaimType {
	return CLAIM_TYPE_ERC721_DEPLOYED
}

// ValidateBasic performs stateless checks
func (e *MsgERC721DeployedClaim) ValidateBasic() error {
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgERC721DeployedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgERC721DeployedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgERC721DeployedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// MsgType should return the action
func (msg MsgERC721DeployedClaim) MsgType() string { return "ERC721_deployed_claim" }

func (msg *MsgERC721DeployedClaim) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// Route should return the name of the module
func (msg MsgERC721DeployedClaim) Route() string { return RouterKey }

// ClaimHash implements BridgeDeposit.Hash
func (msg *MsgERC721DeployedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%s/", msg.CosmosDenom, msg.TokenContract, msg.Name, msg.Symbol)
	return tmhash.Sum([]byte(path))
}

func (msg MsgERC721DeployedClaim) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

func (m *MsgERC721DeployedClaim) GetEventNonce() uint64 {
	if m != nil {
		return m.EventNonce
	}
	return 0
}

func (m *MsgERC721DeployedClaim) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

// GetType returns the type of the claim
func (e *MsgLogicCallExecutedClaim) GetType() ClaimType {
	return CLAIM_TYPE_LOGIC_CALL_EXECUTED
}

// ValidateBasic performs stateless checks
func (e *MsgLogicCallExecutedClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgLogicCallExecutedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

func (msg MsgLogicCallExecutedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgLogicCallExecutedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// MsgType should return the action
func (msg MsgLogicCallExecutedClaim) MsgType() string { return "Logic_Call_Executed_Claim" }

// Route should return the name of the module
func (msg MsgLogicCallExecutedClaim) Route() string { return RouterKey }

func (msg *MsgLogicCallExecutedClaim) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// ClaimHash implements BridgeDeposit.Hash
func (msg *MsgLogicCallExecutedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", msg.InvalidationId, msg.InvalidationNonce)
	return tmhash.Sum([]byte(path))
}

func (msg MsgLogicCallExecutedClaim) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return acc
}

// NewMsgCancelSendToEth returns a new msgSetOrchestratorAddress
func NewMsgCancelSendToEth(val sdk.AccAddress, id uint64) *MsgCancelSendToEth {
	return &MsgCancelSendToEth{
		TransactionId: id,
	}
}

// Route should return the name of the module
func (msg *MsgCancelSendToEth) Route() string { return RouterKey }

// MsgType should return the action
func (msg *MsgCancelSendToEth) MsgType() string { return "cancel_send_to_eth" }

// ValidateBasic performs stateless checks
func (msg *MsgCancelSendToEth) ValidateBasic() (err error) {
	_, err = sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	return nil
}

func (msg *MsgCancelSendToEth) Bytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelSendToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(GravityCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgCancelSendToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

func (msg MsgCancelSendToEth) GetFromAddress() sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return acc
}
