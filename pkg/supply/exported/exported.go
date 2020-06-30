package exported

import (
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	)

type ModuleAccountI interface {
	exported.Account

	GetName() string
	GetPermissions() []string
	HasPermission(string) bool
}


// SupplyI defines an inflationary supply interface for modules that handle
// token supply.
type SupplyI interface {
	GetTotal() sdk.Coin
	SetTotal(total sdk.Coin) SupplyI

	Inflate(amount sdk.Coin) SupplyI
	Deflate(amount sdk.Coin) SupplyI

	String() string
	ValidateBasic() error
}