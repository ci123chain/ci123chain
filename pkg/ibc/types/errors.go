package types


import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)


type CodeType = sdk.CodeType

// ibc errors reserve 300 ~ 399.
const (
	DefaultCodespace 			sdk.CodespaceType = "ibc"
	CodeBadBankSignature       	CodeType = 301
	CodeBadReceiptSignature		CodeType = 302
	CodeBadUnmarshal      		CodeType = 303
	CodeBadMarshal      		CodeType = 304
	CodeGetBankAddrError		CodeType = 305
	CodeMakeIBCMsgError			CodeType = 306
	CodeSetIBCMsgError			CodeType = 307
	CodeApplyIBCMsgError		CodeType = 308
	CodeMakeBankReceiptError	CodeType = 309
	CodeBankSendError			CodeType = 310
	CodeReceiveReceiptError		CodeType = 311
	CodeBadState				CodeType = 312
)

func ErrBadBankSignature(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeBadBankSignature, "Bank msg verify failed: %s", err.Error())
}

func ErrBadReceiptSignature(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeBadReceiptSignature, "Receipt msg verify failed: %s", err.Error())
}

func ErrFailedUnmarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
	return sdk.NewError(codespace, CodeBadUnmarshal, "Unmarshal Error: %s", detailStr)
}

func ErrFailedMarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
	return sdk.NewError(codespace, CodeBadMarshal, "Marshal Error: %s", detailStr)
}

func ErrGetBankAddr(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGetBankAddrError, "Get bank address error: %s", err.Error())
}

func ErrMakeIBCMsg(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeMakeIBCMsgError, "Make IBCMsg error: %s", err.Error())
}

func ErrSetIBCMsg(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeSetIBCMsgError, "Set IBCMsg error: %s", err.Error())
}

func ErrApplyIBCMsg(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeApplyIBCMsgError, "Apply IBCMsg error: %s", err.Error())
}

func ErrMakeBankReceipt(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeMakeBankReceiptError, "MakeBankReceipt error: %s", err.Error())
}

func ErrReceiveReceipt(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeReceiveReceiptError, "Receive receipt error: %s", err.Error())
}

func ErrBankSend(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeBankSendError, "Bank send error: %s", err.Error())
}

func ErrState(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeBadState, "State error: %s", err.Error())
}