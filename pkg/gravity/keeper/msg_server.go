package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	stakingtypes "github.com/ci123chain/ci123chain/pkg/staking/types"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

type msgServer struct {
	*Keeper
}

const (
	ERC20 = iota
	ERC721
)

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: &keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetGravityID(c context.Context,gid string) error {
	ctx := sdk.UnwrapSDKContext(c)
	if gid == "" {
		return types.ErrEmpty.Wrap(" gravityid")
	}
	k.Keeper.saveGravityID(ctx, gid)
	k.SetCurrentGid(gid)
	return nil
}

func (k msgServer) RevertGravityID() {
	k.SetCurrentGid("")
}

// ValsetConfirm handles MsgValsetConfirm
// TODO: check msgValsetConfirm to have an Orchestrator field instead of a Validator field
func (k msgServer) ValsetConfirm(c context.Context, msg *types.MsgValsetConfirm) (*types.MsgValsetConfirmResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	valset := k.GetValset(ctx, msg.Nonce)
	if valset == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find valset")
	}
	
	gravityID := k.currentGID
	checkpoint := valset.GetCheckpoint(gravityID)
	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}
	// ensure that the validator exists
	orchaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	if k.StakingKeeper.Validator(ctx, orchaddr) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, orchaddr.String())
	}

	if err = types.ValidateEthereumSignature(checkpoint, sigBytes, orchaddr.String()); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with gravity-id %s with checkpoint %s found %s", orchaddr.String(), gravityID, hex.EncodeToString(checkpoint), msg.Signature))
	}

	// persist signature
	if k.GetValsetConfirmByGID(ctx, msg.Nonce, orchaddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature duplicate")
	}
	key := k.SetValsetConfirmByGID(ctx, *msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			sdk.NewAttribute([]byte(types.AttributeKeyValsetConfirmKey), key),
		),
	)

	return &types.MsgValsetConfirmResponse{}, nil
}

// SendToEth handles MsgSendToEth
func (k msgServer) SendToEth(c context.Context, msg *types.MsgSendToEth) (*types.MsgSendToEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	// fee at least 100 wlk
	if msg.BridgeFee.Amount.LT(sdk.NewIntWithDecimal(1, 20)) {
		return nil, errors.New("Bridge fee is less than 100 wlk")
	}

	txID, err := k.AddToOutgoingPool(ctx, sender, msg.EthDest, msg.Amount, msg.BridgeFee, msg.TokenType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			sdk.NewAttribute([]byte(types.AttributeKeyOutgoingTXID), []byte(fmt.Sprint(txID))),
		),
	)

	return &types.MsgSendToEthResponse{}, nil
}

// RequestBatch handles MsgRequestBatch
func (k msgServer) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Check if the denom is a gravity coin, if not, check if there is a deployed ERC20 representing it.
	// If not, error out
	var tokenContract string
	if msg.TokenType == ERC20 {
		_, token, err := k.DenomToERC20Lookup(ctx, msg.Denom)
		if err != nil {
			return nil, err
		}
		tokenContract = token
	} else if msg.TokenType == ERC721 {
		_, token, err := k.DenomToERC721Lookup(ctx, msg.Denom)
		if err != nil {
			return nil, err
		}
		tokenContract = token
	}

	batchID, err := k.BuildOutgoingTXBatch(ctx, tokenContract, OutgoingTxBatchSize, msg.TokenType, msg.GetFromAddress().Address)
	if err != nil {
		return nil, err
	}
	if batchID == nil {
		return nil, types.ErrBatchIDNil
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			sdk.NewAttribute([]byte(types.AttributeKeyBatchNonce), []byte(fmt.Sprint(batchID.BatchNonce))),
		),
	)

	return &types.MsgRequestBatchResponse{}, nil
}

// ConfirmBatch handles MsgConfirmBatch
func (k msgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := k.GetOutgoingTXBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find batch")
	}

	gravityID := k.currentGID
	checkpoint, err := batch.GetCheckpoint(gravityID, msg.TokenType)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	}

	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	orchaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	if k.StakingKeeper.Validator(ctx, orchaddr) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, orchaddr.String())
	}
	err = types.ValidateEthereumSignature(checkpoint, sigBytes, orchaddr.String())
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with gravity-id %s with checkpoint %s found %s", orchaddr.String(), gravityID, hex.EncodeToString(checkpoint), msg.Signature))
	}

	// check if we already have this confirm
	if k.GetBatchConfirmWithGID(ctx, msg.Nonce, msg.TokenContract, orchaddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
	}
	key := k.SetBatchConfirmWithGID(ctx, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			sdk.NewAttribute([]byte(types.AttributeKeyBatchConfirmKey), key),
		),
	)

	return nil, nil
}

// ConfirmLogicCall handles MsgConfirmLogicCall
//func (k msgServer) ConfirmLogicCall(c context.Context, msg *types.MsgConfirmLogicCall) (*types.MsgConfirmLogicCallResponse, error) {
//	ctx := sdk.UnwrapSDKContext(c)
//	invalidationIdBytes, err := hex.DecodeString(msg.InvalidationId)
//	if err != nil {
//		return nil, sdkerrors.Wrap(types.ErrInvalid, "invalidation id encoding")
//	}
//
//	// fetch the outgoing logic given the nonce
//	logic := k.GetOutgoingLogicCall(ctx, invalidationIdBytes, msg.InvalidationNonce)
//	if logic == nil {
//		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find logic")
//	}
//
//	gravityID := k.currentGID
//	checkpoint, err := logic.GetCheckpoint(gravityID)
//	if err != nil {
//		return nil, sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
//	}
//
//	sigBytes, err := hex.DecodeString(msg.Signature)
//	if err != nil {
//		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
//	}
//
//	orchaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
//	if k.StakingKeeper.Validator(ctx, orchaddr) == nil {
//		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, orchaddr.String())
//	}
//
//	ethAddress := orchaddr.String()
//
//	err = types.ValidateEthereumSignature(checkpoint, sigBytes, ethAddress)
//	if err != nil {
//		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with gravity-id %s with checkpoint %s found %s", ethAddress, gravityID, hex.EncodeToString(checkpoint), msg.Signature))
//	}
//
//	// check if we already have this confirm
//	if k.GetLogicCallConfirm(ctx, invalidationIdBytes, msg.InvalidationNonce, orchaddr) != nil {
//		return nil, sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
//	}
//
//	k.SetLogicCallConfirm(ctx, msg)
//
//	ctx.EventManager().EmitEvent(
//		sdk.NewEvent(
//			sdk.EventTypeMessage,
//			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
//		),
//	)
//
//	return nil, nil
//}

// DepositClaim handles MsgDepositClaim
// TODO it is possible to submit an old msgDepositClaim (old defined as covering an event nonce that has already been
// executed aka 'observed' and had it's slashing window expire) that will never be cleaned up in the endblocker. This
// should not be a security risk as 'old' events can never execute but it does store spam in the chain.
func (k msgServer) DepositClaim(c context.Context, msg *types.MsgDepositClaim) (*types.MsgDepositClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	orchaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)

	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, orchaddr)
	if val == nil || !val.IsBonded() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in active set")
	}

	//any, err := codectypes.NewAnyWithValue(msg)
	//if err != nil {
	//	return nil, err
	//}

	// Add the claim to the store
	_, err := k.Attest(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute([]byte(types.AttributeKeyAttestationID), types.GetAttestationKey(msg.EventNonce, msg.ClaimHash())),
		),
	)

	k.setEventNonceState(ctx, msg.EventNonce, eventNonceStateDone)

	return &types.MsgDepositClaimResponse{}, nil
}

// WithdrawClaim handles MsgWithdrawClaim
// TODO it is possible to submit an old msgWithdrawClaim (old defined as covering an event nonce that has already been
// executed aka 'observed' and had it's slashing window expire) that will never be cleaned up in the endblocker. This
// should not be a security risk as 'old' events can never execute but it does store spam in the chain.
func (k msgServer) WithdrawClaim(c context.Context, msg *types.MsgWithdrawClaim) (*types.MsgWithdrawClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	orchaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)

	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, orchaddr)
	if val == nil || !val.IsBonded() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
	}

	// Add the claim to the store
	_, err := k.Attest(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute([]byte(types.AttributeKeyAttestationID), types.GetAttestationKey(msg.EventNonce, msg.ClaimHash())),
		),
	)

	return &types.MsgWithdrawClaimResponse{}, nil
}

// ERC20Deployed handles MsgERC20Deployed
func (k msgServer) ERC20DeployedClaim(c context.Context, msg *types.MsgERC20DeployedClaim) (*types.MsgERC20DeployedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orchAddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)

	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, orchAddr)
	if val == nil || !val.IsBonded() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
	}

	// Add the claim to the store
	_, err := k.Attest(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute([]byte(types.AttributeKeyAttestationID), types.GetAttestationKey(msg.EventNonce, msg.ClaimHash())),
		),
	)

	return &types.MsgERC20DeployedClaimResponse{}, nil
}


func (k msgServer) CancelSendToEth(c context.Context, msg *types.MsgCancelSendToEth) (*types.MsgCancelSendToEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender := sdk.HexToAddress(msg.Sender)
	err := k.RemoveFromOutgoingPoolAndRefund(ctx, msg.TransactionId, sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			sdk.NewAttribute([]byte(types.AttributeKeyOutgoingTXID), []byte(fmt.Sprint(msg.TransactionId))),
		),
	)

	return &types.MsgCancelSendToEthResponse{}, nil
}

func (k msgServer) Deposit721Claim(c context.Context, msg *types.MsgDeposit721Claim) (*types.MsgDeposit721ClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	orchaddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, orchaddr)
	if val == nil || !val.IsBonded() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in active set")
	}

	// Add the claim to the store
	_, err := k.Attest(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute([]byte(types.AttributeKeyAttestationID), types.GetAttestationKey(msg.EventNonce, msg.ClaimHash())),
		),
	)

	k.setEventNonceState(ctx, msg.EventNonce, eventNonceStateDone)

	return &types.MsgDeposit721ClaimResponse{}, nil
}

func (k msgServer) ERC721DeployedClaim(c context.Context, msg *types.MsgERC721DeployedClaim) (*types.MsgERC721DeployedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orchAddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, orchAddr)
	if val == nil || !val.IsBonded() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
	}

	// Add the claim to the store
	_, err := k.Attest(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute([]byte(types.AttributeKeyAttestationID), types.GetAttestationKey(msg.EventNonce, msg.ClaimHash())),
		),
	)

	return &types.MsgERC721DeployedClaimResponse{}, nil
}

func (k msgServer) ValsetConfirmNonceClaim(c context.Context, msg *types.MsgValsetConfirmNonceClaim) (*types.MsgValsetConfirmNonceClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orchAddr, _ := sdk.AccAddressFromBech32(msg.Orchestrator)

	// return an error if the validator isn't in the active set
	val := k.StakingKeeper.Validator(ctx, orchAddr)
	if val == nil || !val.IsBonded() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, "validator not in acitve set")
	}

	// Add the claim to the store
	_, err := k.Attest(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation error")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(msg.MsgType())),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute([]byte(types.AttributeKeyAttestationID), types.GetAttestationKey(msg.ValsetNonce, msg.ClaimHash())),
		),
	)

	return &types.MsgValsetConfirmNonceClaimResponse{}, nil
}
