package gravity

import (
	"fmt"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"

	"github.com/ci123chain/ci123chain/pkg/gravity/keeper"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)



// NewHandler returns a handler for "Gravity" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		if gMsg, ok := msg.(types.GravityInterface); !ok {
			return sdk.WrapServiceResult(ctx, nil, types.ErrInvalid)
		} else {
			gravityID := gMsg.GetGravityID()
			if err := msgServer.SetGravityID(sdk.WrapSDKContext(ctx), gravityID); err != nil {
				return sdk.WrapServiceResult(ctx, nil, err)
			}
		}

		var res interface{}
		var err error
		switch msg := msg.(type) {
		case *types.MsgValsetConfirm:
			res, err = msgServer.ValsetConfirm(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgSendToEth:
			res, err = msgServer.SendToEth(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgCancelSendToEth:
			res, err = msgServer.CancelSendToEth(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgRequestBatch:
			res, err = msgServer.RequestBatch(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgConfirmBatch:
			res, err = msgServer.ConfirmBatch(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgDepositClaim:
			res, err = msgServer.DepositClaim(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgDeposit721Claim:
			res, err = msgServer.Deposit721Claim(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgWithdrawClaim:
			res, err = msgServer.WithdrawClaim(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgERC20DeployedClaim:
			res, err = msgServer.ERC20DeployedClaim(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgERC721DeployedClaim:
			res, err = msgServer.ERC721DeployedClaim(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgValsetConfirmNonceClaim:
			res, err = msgServer.ValsetConfirmNonceClaim(sdk.WrapSDKContext(ctx), msg)
		default:
			res, err = nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Gravity Msg type: %v", msg.MsgType()))
		}

		msgServer.RevertGravityID()
		return sdk.WrapServiceResult(ctx, res, err)
	}
}
