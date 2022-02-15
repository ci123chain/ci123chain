package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/account"
	supplytypes "github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
)

var MAX_UINT, _ = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457", 10)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	keeper     Keeper
	supplyKeeper supply.Keeper
	accountKeeper account.AccountKeeper
}

// Handle is the entry point for Attestation processing.
// TODO-JT add handler for ERC20DeployedEvent
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.EthereumClaim) error {
	ma := a.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	defer func(account supplytypes.ModuleAccountI) {
		if err := account.SetSequence(account.GetSequence() + 1); err != nil {
			panic(err)
		}
		a.accountKeeper.SetAccount(ctx, account)
	}(ma)
	switch claim := claim.(type) {
	case *types.MsgDepositClaim:
		// Check if coin is Cosmos-originated asset and get denom
		exists, wlkToken := a.keeper.ERC20ToDenomLookup(ctx, claim.TokenContract)

		if !exists {
			tokenAddres, err := a.supplyKeeper.DeployWRC20ForGivenERC20(ctx, types.ModuleName, []interface{}{claim.TokenName, claim.TokenSymbol, uint8(claim.TokenDecimals), MAX_UINT, false})
			if err != nil {
				return sdkerrors.Wrap(err, "deploy wrc20 failed")
			}
			wlkToken = tokenAddres.String()
			a.keeper.setERC20Map(ctx, wlkToken, claim.TokenContract)
			a.keeper.SetTokenMetaData(ctx, wlkToken, types.MetaData{
				Symbol:   claim.TokenSymbol,
				Name:     claim.TokenName,
				Decimals: claim.TokenDecimals,
			})
		}
		var err error
		if a.keeper.IsWlkToken(wlkToken) {
			// 如果是 weelink 原生币
			// If it is cosmos originated, unlock the coins
			coins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, claim.Amount)}

			addr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid reciever address")
			}

			if err = a.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
				err = sdkerrors.Wrap(err, "transfer vouchers")
			}
		} else {
			// erc20 代币
			receiverAddr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
			if err != nil {
				return sdkerrors.Wrap(err, "invalid reciever address")
			}

			metaData := a.keeper.GetTokenMetaData(ctx, wlkToken)
			if metaData.Symbol == "" {
				return types.ErrNoContractMetaData
			}

			//err = a.supplyKeeper.MintCoinsFromModuleToEvmAccount(ctx, receiverAddr, wlkToken, claim.Amount.BigInt())
			err = a.supplyKeeper.TransferFromModuleToEvmAccount(ctx, receiverAddr, wlkToken, claim.Amount.BigInt())
		}
		return err
	case *types.MsgDeposit721Claim:
		// Check if coin is Cosmos-originated asset and get denom
		exists, wlkToken := a.keeper.ERC721ToDenomLookup(ctx, claim.TokenContract)

		if !exists {
			tokenAddres, err := a.supplyKeeper.DeployWRC721ForGivenERC721(ctx, types.ModuleName, []interface{}{claim.TokenName, claim.TokenSymbol})
			if err != nil {
				return sdkerrors.Wrap(err, "deploy wrc20 failed")
			}
			wlkToken = tokenAddres.String()
			a.keeper.setERC721Map(ctx, wlkToken, claim.TokenContract)
			a.keeper.SetTokenMetaData(ctx, wlkToken, types.MetaData{
				Symbol:   claim.TokenSymbol,
				Name:     claim.TokenName,
			})
		}
		var err error
		// erc20 代币
		receiverAddr, err := sdk.AccAddressFromBech32(claim.CosmosReceiver)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid reciever address")
		}

		metaData := a.keeper.GetTokenMetaData(ctx, wlkToken)
		if metaData.Symbol == "" {
			return types.ErrNoContractMetaData
		}

		//err = a.supplyKeeper.MintCoinsFromModuleToEvmAccount(ctx, receiverAddr, wlkToken, claim.Amount.BigInt())
		err = a.supplyKeeper.Transfer721FromModuleToEvmAccount(ctx, receiverAddr, wlkToken, claim.TokenID.BigInt())
		return err
	case *types.MsgWithdrawClaim:
		a.keeper.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)
	case *types.MsgERC20DeployedClaim:
		wlkToken := claim.CosmosDenom
		// Check if it already exists
		existERC20, exists := a.keeper.GetMapedEthToken(ctx, wlkToken)
		if exists {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("WRC20 %s in weelink already exists for eth %s", wlkToken, existERC20))
		}

		existingWRC20, exists := a.keeper.GetMapedWlkToken(ctx, claim.TokenContract)
		if exists {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC20 %s in eth already exists for weelink %s", claim.TokenContract, existingWRC20))
		}

		if err := a.validateAndSetWRC20(ctx, wlkToken, claim.Name, claim.Symbol, claim.Decimals); err != nil {
			return err
		}
		// Add to wrc20-erc20 mapping
		a.keeper.setERC20Map(ctx, wlkToken, claim.TokenContract)
		a.logger(ctx).Info("ERC20-WRC20 Mapped for ", "ETH:", claim.TokenContract, " Weelink:", wlkToken)
	case *types.MsgERC721DeployedClaim:
		wlkToken := claim.CosmosDenom
		// Check if it already exists
		existERC721, exists := a.keeper.GetMapedERC721Token(ctx, wlkToken)
		if exists {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("WRC721 %s in weelink already exists for eth %s", wlkToken, existERC721))
		}

		existingWRC721, exists := a.keeper.GetMapedWRC721Token(ctx, claim.TokenContract)
		if exists {
			return sdkerrors.Wrap(
				types.ErrInvalid,
				fmt.Sprintf("ERC721 %s in eth already exists for weelink %s", claim.TokenContract, existingWRC721))
		}

		if err := a.validateAndSetWRC721(ctx, wlkToken, claim.Name, claim.Symbol); err != nil {
			return err
		}
		// Add to wrc20-erc20 mapping
		a.keeper.setERC721Map(ctx, wlkToken, claim.TokenContract)
		a.logger(ctx).Info("ERC721-WRC721 Mapped for ", "ETH:", claim.TokenContract, " Weelink:", wlkToken)
	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}

func (a AttestationHandler) validateAndSetWRC20(ctx sdk.Context, denom, name, symbol string, decimals uint64) error {
	// todo validate denom from erc20
	if a.keeper.IsWlkToken(denom){
		// chain coin
		if decimals != 0 {
			return types.ErrDenomDecimal.Wrapf("decimals for weelink expect zero, got %d", decimals)
		}
	} else {
		//
		if err := types.ValidateEthAddress(denom); err != nil {
			return err
		}
		name2, err := a.supplyKeeper.WRC20DenomValueForFunc(ctx, types.ModuleName, sdk.HexToAddress(denom), "name")
		if err != nil {
			return types.ErrQueryDenom.Wrap(err.Error())
		}
		symbol2, err := a.supplyKeeper.WRC20DenomValueForFunc(ctx, types.ModuleName, sdk.HexToAddress(denom), "symbol")
		if err != nil {
			return types.ErrQueryDenom.Wrap(err.Error())
		}
		decimal2, err := a.supplyKeeper.WRC20DenomValueForFunc(ctx, types.ModuleName, sdk.HexToAddress(denom), "decimals")
		if err != nil {
			return types.ErrQueryDenom.Wrap(err.Error())
		}
		nameS, ok := name2.(string)
		if !ok {
			return types.ErrQueryDenom.Wrap("invalid name type expected string")
		}
		symbolS, ok := symbol2.(string)
		if !ok {
			return types.ErrQueryDenom.Wrap("invalid symbol type expected string")
		}
		decimalI, ok := decimal2.(uint8)
		if !ok {
			return types.ErrQueryDenom.Wrap("invalid decimal type expected uint64")
		}
		if nameS != name ||
			symbolS != symbol ||
			decimalI != uint8(decimals) {
			return types.ErrQueryDenomMismatch.Wrapf("expect name:%s, symbol:%s, decimals:%d, got name:%s, symbol:%s, decimals:%d", name2, symbol2, decimal2, name, symbol, decimals)
		}
	}

	a.keeper.SetTokenMetaData(ctx, denom, types.MetaData{
		Symbol:   symbol,
		Name:     name,
		Decimals: decimals,
	})
	return nil
}

func (a AttestationHandler) validateAndSetWRC721(ctx sdk.Context, denom, name, symbol string) error {
	if err := types.ValidateEthAddress(denom); err != nil {
		return err
	}
	name2, err := a.supplyKeeper.WRC20DenomValueForFunc(ctx, types.ModuleName, sdk.HexToAddress(denom), "name")
	if err != nil {
		return types.ErrQueryDenom.Wrap(err.Error())
	}
	symbol2, err := a.supplyKeeper.WRC20DenomValueForFunc(ctx, types.ModuleName, sdk.HexToAddress(denom), "symbol")
	if err != nil {
		return types.ErrQueryDenom.Wrap(err.Error())
	}
	nameS, ok := name2.(string)
	if !ok {
		return types.ErrQueryDenom.Wrap("invalid name type expected string")
	}
	symbolS, ok := symbol2.(string)
	if !ok {
		return types.ErrQueryDenom.Wrap("invalid symbol type expected string")
	}
	if nameS != name ||
		symbolS != symbol {
		return types.ErrQueryDenomMismatch.Wrapf("expect name:%s, symbol:%s, got name:%s, symbol:%s", name2, symbol2, name, symbol)
	}

	a.keeper.SetTokenMetaData(ctx, denom, types.MetaData{
		Symbol:   symbol,
		Name:     name,
	})
	return nil
}

// logger returns a module-specific logger.
func (k AttestationHandler) logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}