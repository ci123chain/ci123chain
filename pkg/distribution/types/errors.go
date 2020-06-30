package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace  sdk.CodespaceType = "distribution"

	CodeInvalidHeight CodeType = 300

	CodeBadMarshal  CodeType = 301

	CodeBadAddress  CodeType = 302

	CodeEmptyDelegationStartingInfo CodeType = 303

	CodeInternalServerError CodeType = 304

	CodeInvalidCoin  CodeType   = 305
	CodeInvalidAddress CodeType  = 306
	CodeInvalidAmount  CodeType  = 307
	CodeInvalidGas    CodeType   = 308
	CodeInvalidPrivateKey CodeType = 309
	CodeErrSignTx   CodeType    = 310
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