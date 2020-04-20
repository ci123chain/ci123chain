package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"strings"
)

// permissions
const (
	Burner  = "burner"
	Staking = "staking"
)

type PermissionsForAddress struct {
	permissions 	[]string
	address 		sdk.AccAddress
}

func NewPermissionForAddress(name string, permissions []string) PermissionsForAddress {
	return PermissionsForAddress{
		permissions: 	permissions,
		address: 		NewModuleAddress(name),
	}
}


// GetAddress returns the address of the PermissionsForAddress object
func (pa PermissionsForAddress) GetAddress() sdk.AccAddress {
	return pa.address
}

// GetPermissions returns the permissions granted to the address
func (pa PermissionsForAddress) GetPermissions() []string {
	return pa.permissions
}

// performs basic permission validation
func validatePermissions(permissions ...string) error {
	for _, perm := range permissions {
		if strings.TrimSpace(perm) == "" {
			return fmt.Errorf("module permission is empty")
		}
	}
	return nil
}
