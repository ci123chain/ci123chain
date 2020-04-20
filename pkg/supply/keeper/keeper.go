package keeper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/supply/exported"
	types2 "github.com/ci123chain/ci123chain/pkg/supply/types"
)

type Keeper struct {
	cdc 		*codec.Codec
	storeKey 	sdk.StoreKey
	ak 			account.AccountKeeper

	permAddrs 	map[string]types2.PermissionsForAddress
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ak account.AccountKeeper, maccPerms map[string][]string) Keeper {
	permAddrs := make(map[string]types2.PermissionsForAddress)
	for name, perms := range maccPerms {
		permAddrs[name] = types2.NewPermissionForAddress(name, perms)
	}

	return Keeper{
		cdc: cdc,
		storeKey: key,
		ak: 	ak,
		permAddrs: permAddrs,
	}
}

func (k Keeper) GetModuleAccount(ctx sdk.Context, moduleName string ) exported.ModuleAccountI {
	acc, _ := k.GetModuleAccountAndPermissions(ctx, moduleName)
	return acc
}

// GetModuleAddress returns an address based on the module name
func (k Keeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return sdk.AccAddress{}
	}
	return permAddr.GetAddress()
}

func (k Keeper) GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (exported.ModuleAccountI, []string) {
	addr, perms := k.GetModuleAddressAndPermissions(moduleName)
	if addr.Empty() {
		return nil, []string{}
	}

	acc := k.ak.GetAccount(ctx, addr)
	if acc != nil {
		macc, ok := acc.(exported.ModuleAccountI)
		if !ok {
			panic("account is not a module account")
		}
		return macc, perms
	}

	macc := types2.NewEmptyModuleAccount(moduleName, perms...)
	maccI := (k.ak.NewAccount(ctx, macc)).(exported.ModuleAccountI)
	k.SetModuleAccount(ctx, maccI)
	return maccI, perms
}


func (k Keeper) GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string) {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return addr, permissions
	}
	return permAddr.GetAddress(), permAddr.GetPermissions()
}

func (k Keeper) SetModuleAccount(ctx sdk.Context, macc exported.ModuleAccountI) {
	k.ak.SetAccount(ctx, macc)
}



// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount
func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress,
	recipientModule string, amt sdk.Coin) sdk.Error {

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(fmt.Sprintf("module account %s isn't able to be created", recipientModule))
	}
	return k.ak.Transfer(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string,
	recipientAddr sdk.AccAddress, amt sdk.Coin) sdk.Error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}


	return k.ak.Transfer(ctx, senderAddr, recipientAddr, amt)
}

func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coin) error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s isn't able to be created", recipientModule))
	}
	return k.ak.Transfer(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

func (k Keeper) DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string,
	amt sdk.Coin) error {

	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s isn't able to be created", recipientModule))
	}

	if !recipientAcc.HasPermission(types2.Staking) {
		return sdk.ErrNoPermission(fmt.Sprintf("module account %s has no expected permission", recipientModule))
	}
	return k.ak.Transfer(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// UndelegateCoinsFromModuleToAccount undelegates the unbonding coins and transfers
// them from a module account to the delegator account. It will panic if the
// module account does not exist or is unauthorized.
func (k Keeper) UndelegateCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coin,
) error {

	acc := k.GetModuleAccount(ctx, senderModule)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s isn't able to be created", recipientAddr))
	}

	if !acc.HasPermission(types2.Staking) {
		return sdk.ErrNoPermission(fmt.Sprintf("module account %s has no expected permission", recipientAddr))
	}

	return k.ak.Transfer(ctx, acc.GetAddress(), recipientAddr, amt)
}
///-------------

//func (k Keeper) SetAccountSequence(ctx sdk.Context, addr sdk.AccAddress, nonce uint64) sdk.Error {
//	k.ak.SetSequence(ctx, addr, nonce)
//	return nil
//}

///-------------
