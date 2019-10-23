package transaction

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

func GetNonceByAddress(add common.Address) (uint64, error) {
	return uint64(time.Now().UnixNano()), nil
}
