package exported

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

//type ModuleAccountI interface {
//	exported.Account
//
//	GetName() string
//	GetPermissions() []string
//	HasPermission(string) bool
//}
// SupplyI defines an inflationary supply interface for modules that handle
// token supply.
type SupplyI interface {
	GetTotal() sdk.Coins
	SetTotal(total sdk.Coins) SupplyI

	Inflate(amount sdk.Coins) SupplyI
	Deflate(amount sdk.Coins) SupplyI

	String() string
	ValidateBasic() error
}