package types


import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)


type CodeType = sdk.CodeType

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace 			sdk.CodespaceType = "ibc"
	CodeBadBankSignature       	CodeType = 301
	CodeBadReceiptSignature		CodeType = 302
	CodeBadPubkey				CodeType = 303
	CodeBadUnmarshal      		CodeType = 305
)

func ErrBadBankSignature(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeBadBankSignature, "Bank msg verify failed")
}

func ErrBadReceiptSignature(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeBadReceiptSignature, "Receipt msg verify failed")
}

func ErrFailedUnmarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
	return sdk.NewError(codespace, CodeBadUnmarshal, "Unmarshal Error: %s", detailStr)
}

func ErrDecodePubkey(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeBadPubkey, "Decode pubkey error")
}