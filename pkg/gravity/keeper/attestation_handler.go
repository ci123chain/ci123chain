package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/supply"
	vmtypes "github.com/ci123chain/ci123chain/pkg/vm/types"

	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper     Keeper
	supplyKeeper supply.Keeper
	accountKeeper account.AccountKeeper
	evmKeeper  vmtypes.Keeper
}

// Handle is the entry point for Attestation processing.
// TODO-JT add handler for ERC20DeployedEvent
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.EthereumClaim) error {
	switch claim := claim.(type) {
	case *types.MsgDepositClaim:
		// Check if coin is Cosmos-originated asset and get denom
		isCosmosOriginated, wlkContract := a.keeper.ERC20ToDenomLookup(ctx, claim.TokenContract)

		if isCosmosOriginated {
			// 如果是 weelink 原生币
			// If it is cosmos originated, unlock the coins
			coins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, claim.Amount)}

			addr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid reciever address")
			}

			if err = a.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
				return sdkerrors.Wrap(err, "transfer vouchers")
			}
		} else {
			// erc20 代币
			// If it is not cosmos originated, mint the coins (aka vouchers)
			//coin := sdk.NewCoin(denom, claim.Amount)
			receiverAddr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid reciever address")
			}
			metaData := a.keeper.GetTokenMetaData(ctx, wlkContract)
			if metaData.Symbol == "" {
				return types.ErrNoContractMetaData
			}

			owner := a.supplyKeeper.GetModuleAddress(types.ModuleName)
			//param := []interface{}{metaData.Name, metaData.Symbol, metaData.Symbol, 0, true}
			//denomAddr, err := a.DeployERC20Contract(ctx, owner, param)
			//if err != nil {
			//	return err
			//}

			err = a.Mint(ctx, sdk.HexToAddress(wlkContract), owner, receiverAddr, claim.Amount.BigInt())
			return err
		}
	case *types.MsgWithdrawClaim:
		a.keeper.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)
	case *types.MsgERC20DeployedClaim:
		// Check if it already exists
		existingERC20, exists := a.keeper.GetMapedEthToken(ctx, claim.TokenContract)
		if exists {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 %s in wlk already exists for eth %s", existingERC20, claim.TokenContract))
		}

		owner := a.supplyKeeper.GetModuleAddress(types.ModuleName)
		param := []interface{}{claim.Symbol, claim.Name, claim.Decimals, 0, true}
		wlkAddr, err := a.DeployERC20Contract(ctx, owner, param)
		if err != nil {
			return err
		}
		// Add to denom-erc20 mapping
		a.keeper.setERC20Map(ctx, wlkAddr.String(), claim.TokenContract)
		a.keeper.SetTokenMetaData(ctx, wlkAddr.String(), types.MetaData{
			Symbol:   claim.Symbol,
			Name:     claim.Name,
			Decimals: claim.Decimals,
		})
		//a.keeper.setCosmosOriginatedDenomToERC20(ctx, claim.CosmosDenom, claim.TokenContract)

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
