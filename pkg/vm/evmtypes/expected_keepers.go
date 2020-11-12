package evmtypes

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	accountexported "github.com/ci123chain/ci123chain/pkg/account/exported"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) accountexported.Account
	GetAllAccounts(ctx sdk.Context) (accounts []accountexported.Account)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) accountexported.Account
	SetAccount(ctx sdk.Context, account accountexported.Account)
	RemoveAccount(ctx sdk.Context, account accountexported.Account)
}
