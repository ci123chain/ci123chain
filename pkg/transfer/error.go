package transfer

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace 	sdk.CodespaceType = "transfer"
	CodeInvalidAmount       CodeType = 201
	CodeInvalidReceiver     CodeType = 202
)

//----------------------------------------
// Error constructors

func ErrBadAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAmount, "param Amount invalid")
}

func ErrBadReceiver(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReceiver, "param To invalid")
}

