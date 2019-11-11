package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
)

func TestUUID(t *testing.T)  {
	uuid := GenerateUniqueID()
	uuid = strings.ReplaceAll(uuid, "-", "")
	fmt.Println(uuid)

	cc, err := hex.DecodeString(uuid)
	if err != nil {
		panic(err)
	}

	aa := hex.EncodeToString(cc)
	fmt.Println(aa)
	assert.Equal(t, aa, uuid)
}