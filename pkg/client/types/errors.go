package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace 				sdk.CodespaceType = "client"
	CodeNewClientCtxError       	CodeType = 601
	CodeGetInputAddrError       	CodeType = 602
	CodeParseAddrError       		CodeType = 603
	CodeNoAddrError       			CodeType = 604
	CodeGetPassPhraseError			CodeType = 605
	CodeGetSignDataError			CodeType = 606
	CodeBroadcastError				CodeType = 607
	CodeGetCheckPasswordError		CodeType = 608
	CodeGetPasswordError			CodeType = 609
	CodePhrasesNotMatchError		CodeType = 610
	CodeNodeError					CodeType = 611
)

func ErrNewClientCtx(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeNewClientCtxError,"New client context Error", err)
}

func ErrGetInputAddrCtx(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGetInputAddrError,"Get input address Error", err)
}

func ErrParseAddr(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeParseAddrError,"Parse address Error", err)
}

func ErrNoAddr(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeNoAddrError,"No address Error", err)
}

func ErrGetPassPhrase(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGetPassPhraseError,"Get pass phrase Error", err)
}

func ErrGetSignData(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGetSignDataError,"Get sign data from tx Error", err)
}

func ErrBroadcast(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeBroadcastError,"Broadcast Error", err)
}

func ErrGetCheckPassword(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGetCheckPasswordError,"Get check password Error", err)
}

func ErrGetPassword(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGetPasswordError,"Get password Error", err)
}

func ErrPhrasesNotMatch(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodePhrasesNotMatchError,"Phrases not match Error", err)
}

func ErrNode(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeNodeError,"Node error", err)
}
