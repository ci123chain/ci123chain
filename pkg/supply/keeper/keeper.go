package keeper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/supply/exported"
	"github.com/ci123chain/ci123chain/pkg/supply/types"
	vmtypes "github.com/ci123chain/ci123chain/pkg/vm/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	cdc 		*codec.Codec
	storeKey 	sdk.StoreKey
	ak 			account.AccountKeeper
	evmKeeper   vmtypes.Keeper

	permAddrs 	map[string]types.PermissionsForAddress
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ak account.AccountKeeper, maccPerms map[string][]string) Keeper {
	permAddrs := make(map[string]types.PermissionsForAddress)
	for name, perms := range maccPerms {
		permAddrs[name] = types.NewPermissionForAddress(name, perms)
	}

	return Keeper{
		cdc: cdc,
		storeKey: key,
		ak: 	ak,
		permAddrs: permAddrs,
	}
}

func (k Keeper) SetVMKeeper(vmkeeper vmtypes.Keeper) Keeper {
	k.evmKeeper = vmkeeper
	return k
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

	macc := types.NewEmptyModuleAccount(moduleName, perms...)
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
	recipientModule string, amt sdk.Coins) error {

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(fmt.Sprintf("module account %s isn't able to be created", recipientModule))
	}
	return k.ak.Transfer(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string,
	recipientAddr sdk.AccAddress, amt sdk.Coins) error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("module account %s does not exist", senderModule))
	}


	return k.ak.Transfer(ctx, senderAddr, recipientAddr, amt)
}

func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coin) error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("module account %s does not exist", senderModule))
	}
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("module account %s isn't able to be created", recipientModule))
	}
	return k.ak.Transfer(ctx, senderAddr, recipientAcc.GetAddress(), sdk.NewCoins(amt))
}

func (k Keeper) DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string,
	amt sdk.Coin) error {

	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("module account %s isn't able to be created", recipientModule))
	}

	if !recipientAcc.HasPermission(types.Staking) {
		return sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("module account %s has no expected permission", recipientModule))
	}
	return k.ak.Transfer(ctx, senderAddr, recipientAcc.GetAddress(), sdk.NewCoins(amt))
}

// UndelegateCoinsFromModuleToAccount undelegates the unbonding coins and transfers
// them from a module account to the delegator account. It will panic if the
// module account does not exist or is unauthorized.
func (k Keeper) UndelegateCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coin,
) error {

	acc := k.GetModuleAccount(ctx, senderModule)
	if acc == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("module account %s isn't able to be created", recipientAddr))
	}

	if !acc.HasPermission(types.Staking) {
		return sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("module account %s has no expected permission", recipientAddr))
	}

	return k.ak.Transfer(ctx, acc.GetAddress(), recipientAddr, sdk.NewCoins(amt))
}
///-------------

//func (k Keeper) SetAccountSequence(ctx sdk.Context, addr sdk.AccAddress, nonce uint64) sdk.Error {
//	k.ak.SetSequence(ctx, addr, nonce)
//	return nil
//}

// GetSupply retrieves the Supply from store
func (k Keeper) GetSupply(ctx sdk.Context) (supply exported.SupplyI) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.SupplyKey)
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

// SetSupply sets the Supply to store
func (k Keeper) SetSupply(ctx sdk.Context, supply exported.SupplyI) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(types.SupplyKey, b)
}


// ValidatePermissions validates that the module account has been granted
// permissions within its set of allowed permissions.
func (k Keeper) ValidatePermissions(macc exported.ModuleAccountI) error {
	permAddr := k.permAddrs[macc.GetName()]
	for _, perm := range macc.GetPermissions() {
		if !permAddr.HasPermission(perm) {
			return fmt.Errorf("invalid module permission %s", perm)
		}
	}
	return nil
}


// MintCoins creates new coins from thin air and adds it to the module account.
// It will panic if the module account does not exist or is unauthorized.
func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		//panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
		panic(fmt.Errorf("module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(types.Minter) {
		//panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to mint tokens", moduleName))
		panic(fmt.Errorf( "module account %s does not have permissions to mint tokens", moduleName))
	}

	_, err := k.ak.AddBalance(ctx, acc.GetAddress(), amt)
	if err != nil {
		return err
	}

	// update total supply
	supply := k.GetSupply(ctx)
	// todo fix inflate coins
	supply = supply.Inflate(amt)

	k.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Debug(fmt.Sprintf("minted %s from %s module account", amt.String(), moduleName))

	return nil
}

func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.ak.Transfer(ctx, fromAddr, toAddr, amt)
}



func (k Keeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(types.Minter) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to mint tokens", moduleName))
	}

	_, err := k.ak.AddBalance(ctx, acc.GetAddress(), amt)
	if err != nil {
		return err
	}

	// update total supply
	supply := k.GetSupply(ctx)
	// todo fix to coins
	supply.Inflate(amt)


	k.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info("burned coins from module account", "amount", amt.String(), "from", moduleName)
	return nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetDenomMetaData retrieves the denomination metadata
func (k Keeper) GetDenomMetaData(ctx sdk.Context, denom string) types.Metadata {
	st := ctx.KVStore(k.storeKey)
	st = store.NewPrefixStore(st, types.DenomMetadataKey(denom))

	bz := st.Get([]byte(denom))
	if bz == nil {
		return types.Metadata{}
	}

	var metadata types.Metadata
	k.cdc.MustUnmarshalBinaryBare(bz, &metadata)

	return metadata
}