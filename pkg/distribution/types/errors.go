package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace  						sdk.CodespaceType = "distribution"

	CodeInvalidHeight 						CodeType = 300

	CodeBadMarshal  						CodeType = 301

	CodeBadAddress  						CodeType = 302

	CodeEmptyDelegationStartingInfo 		CodeType = 303

	CodeInternalServerError 				CodeType = 304

	CodeInvalidCoin  						CodeType   = 305
	CodeInvalidAddress 						CodeType  = 306
	CodeInvalidAmount  						CodeType  = 307
	CodeInvalidGas    						CodeType   = 308
	CodeInvalidPrivateKey 					CodeType = 309
	CodeErrSignTx   						CodeType    = 310
	CodeErrSetWithdrawAddrDisabled  		CodeType = 311
	CodeErrNoValidatorDistInfo      		CodeType = 312
	CodeErrEmptyDelegationDistInfo          CodeType = 313
	CodeErrNoValidatorCommission            CodeType = 314
	CodeErrWithdrawAddressInfoMismatch      CodeType = 315
	CodeErrErrHandleTxFailed                CodeType = 316
	CodeErrInvalidSignature                 CodeType = 317
)


func ErrBadHeight(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHeight, "param Height invalid: %s", err.Error())
}

func ErrFailedMarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
	return sdk.NewError(codespace, CodeBadMarshal, "Marshal Error: %s", detailStr)
}

func ErrBadAddress(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeBadAddress, "param address invalid: %s", err.Error())
}

func ErrEmptyDelegationStartingInfo(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyDelegationStartingInfo, "empty delegation starting info")
}

func ErrNoValidatorExist(codespace sdk.CodespaceType, valAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyDelegationStartingInfo, fmt.Sprintf("validator %s not exist", valAddr))
}

func ErrNoDelegationExist(codespace sdk.CodespaceType, valAddr, delAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyDelegationStartingInfo, fmt.Sprintf("Delegatin between %s and %s not exist", valAddr, delAddr))
}

func ErrInternalServer(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInternalServerError, fmt.Sprintf("got information failed"))
}

func ErrInvalidCoin(codespace sdk.CodespaceType, coin string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidCoin, fmt.Sprintf("invalid coin %s", coin))
}

func ErrInvalidAddress(codespace sdk.CodespaceType, address string) sdk.Error{
	return sdk.NewError(codespace, CodeInvalidAddress, fmt.Sprintf("invalid address %s", address))
}

func ErrBasAmount(codespace sdk.CodespaceType, amount string) sdk.Error{
	return sdk.NewError(codespace, CodeInvalidAmount, fmt.Sprintf("invalid amount %s", amount))
}

func ErrGas(codespace sdk.CodespaceType, gas string) sdk.Error{
	return sdk.NewError(codespace, CodeInvalidGas, fmt.Sprintf("invalid gas %s", gas))
}

func ErrEmptyPrivateKey(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPrivateKey, "private key can not be empty")
}

func ErrSignTx(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeErrSignTx, err.Error())
}

func ErrSetWithdrawAddrDisabled(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeErrSetWithdrawAddrDisabled, "set withdraw address disabled")
}

func ErrNoValidatorDistInfo(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeErrNoValidatorDistInfo, "no validator distribution info")
}

func ErrEmptyDelegationDistInfo(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeErrEmptyDelegationDistInfo, "no delegation distribution info")
}

func ErrNoValidatorCommission(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeErrNoValidatorCommission, "no validator commission")
}

func ErrWithdrawAddressInfoMismatch(codespace sdk.CodespaceType, expectedAddr, gotAddr sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeErrWithdrawAddressInfoMismatch, fmt.Sprintf("account address mismatch, expected %s, got %s", expectedAddr.String(), gotAddr.String()))
}

func ErrHandleTxFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeErrErrHandleTxFailed, err.Error())
}

func ErrInvalidSignature(codeSpace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codeSpace, CodeErrInvalidSignature, msg)
}