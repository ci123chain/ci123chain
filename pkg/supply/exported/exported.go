package exported

import "github.com/ci123chain/ci123chain/pkg/account/exported"

type ModuleAccountI interface {
	exported.Account

	GetName() string
	GetPermission() []string
	HasPermission(string) bool
}
