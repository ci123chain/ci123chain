package keeper

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func GenerateUniqueID(b []byte) string {
	hSum := md5.Sum([]byte(b))
	hexString := hex.EncodeToString(hSum[:])
	return strings.ToUpper(hexString)
}
