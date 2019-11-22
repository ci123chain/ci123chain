package nonce

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

func CheckIBCNonce(ctx sdk.Context, savedNonce, nonce uint64) bool {
	//savedSequence := k.AccountKeeper.GetAccount(ctx, addr).GetSequence()
	//sequence := nonce
	if savedNonce != nonce {
		return false
	}
	return true
}

func CheckTransferNonce(ctx sdk.Context, savedSequence, nonce uint64) bool {
	//savedSequence := k.GetAccount(ctx, addr).GetSequence()
	sequence := nonce
	if savedSequence != sequence {
		return false
	}
	return true
}

