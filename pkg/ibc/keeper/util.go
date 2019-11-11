package keeper

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GenerateUniqueID() string {
	u1 := uuid.NewV4()
	//fmt.Printf("UUIDv4: %s\n", u1)
	return strings.ReplaceAll(u1.String(), "-", "")
}
