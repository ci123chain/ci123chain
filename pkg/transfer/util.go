package transfer

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"time"
)

func GetNonceByAddress(add types.AccAddress) (uint64, error) {


	return uint64(time.Now().UnixNano()), nil
}
