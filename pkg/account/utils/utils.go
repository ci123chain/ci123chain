package utils

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
)

func QueryNonce(ctx context.Context, addr sdk.AccAddress) (uint64, error) {
	nonce, _, err := ctx.GetNonceByAddress(addr, false)
	return nonce, err
}