package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace  						sdk.CodespaceType = "distribution"

	//CodeInvalidHeight 						CodeType = 300
	//
	//CodeBadMarshal  						CodeType = 301

	//CodeBadAddress  						CodeType = 302

	CodeEmptyDelegationStartingInfo 		CodeType = 1403

	CodeNoValidatorExist                    CodeType = 1404

	CodeCdcMarshalFailed                    CodeType = 1405

	CodeNoDelegationExist                   CodeType = 1406

	CodeSetWithdrawAddressFailed            CodeType = 1407

	CodeNoValidatorDistInfo                 CodeType = 1408

	CodeNoDelegationDistInfo                CodeType = 1409

	CodeNoValidatorCommission               CodeType = 1410

	//CodeInternalServerError 				CodeType = 304

	//CodeInvalidCoin  						CodeType   = 305
	//CodeInvalidAddress 						CodeType  = 306
	//CodeInvalidAmount  						CodeType  = 307
	//CodeInvalidGas    						CodeType   = 308
	//CodeInvalidPrivateKey 					CodeType = 309
	//CodeErrSignTx   						CodeType    = 310
	//CodeErrSetWithdrawAddrDisabled  		CodeType = 311
	//CodeErrNoValidatorDistInfo      		CodeType = 312
	//CodeErrEmptyDelegationDistInfo          CodeType = 313
	//CodeErrNoValidatorCommission            CodeType = 314
	//CodeErrWithdrawAddressInfoMismatch      CodeType = 315
	//CodeErrErrHandleTxFailed                CodeType = 316
	//CodeErrInvalidSignature                 CodeType = 317

	CodeErrInvalidParams                    CodeType = 1411
)


func ErrInvalidParams(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeErrInvalidParams), desc)
}


//func ErrBadHeight(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeInvalidHeight, "param Height invalid: %s", err.Error())
//}
//
//func ErrFailedMarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
//	return sdk.NewError(codespace, CodeBadMarshal, "Marshal Error: %s", detailStr)
//}

//func ErrBadAddress(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeBadAddress, "param address invalid: %s", err.Error())
//}

func ErrEmptyDelegationStartingInfo(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeEmptyDelegationStartingInfo), desc)
}

func ErrNoValidatorExist(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeNoValidatorExist), desc)
}

func ErrInternalCdcMarshal(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcMarshalFailed), desc)
}

func ErrNoDelegationExist(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeNoDelegationExist), desc)
}

//func ErrInternalServer(codespace sdk.CodespaceType) sdk.Error {
//	return sdk.NewError(codespace, CodeInternalServerError, fmt.Sprintf("got information failed"))
//}

//func ErrInvalidCoin(codespace sdk.CodespaceType, coin string) sdk.Error {
//	return sdk.NewError(codespace, CodeInvalidCoin, fmt.Sprintf("invalid coin %s", coin))
//}

//func ErrInvalidAddress(codespace sdk.CodespaceType, address string) sdk.Error{
//	return sdk.NewError(codespace, CodeInvalidAddress, fmt.Sprintf("invalid address %s", address))
//}

//func ErrBasAmount(codespace sdk.CodespaceType, amount string) sdk.Error{
//	return sdk.NewError(codespace, CodeInvalidAmount, fmt.Sprintf("invalid amount %s", amount))
//}
//
//func ErrGas(codespace sdk.CodespaceType, gas string) sdk.Error{
//	return sdk.NewError(codespace, CodeInvalidGas, fmt.Sprintf("invalid gas %s", gas))
//}
//
//func ErrEmptyPrivateKey(codespace sdk.CodespaceType) sdk.Error {
//	return sdk.NewError(codespace, CodeInvalidPrivateKey, "private key can not be empty")
//}
//
//func ErrSignTx(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeErrSignTx, err.Error())
//}

//func ErrParams(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeErrSignTx, err.Error())
//}

func ErrSetWithdrawAddrDisabled(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeSetWithdrawAddressFailed), desc)
}

func ErrNoValidatorDistInfo(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeNoValidatorDistInfo), desc)
}

func ErrEmptyDelegationDistInfo(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeNoDelegationDistInfo), desc)
}

func ErrNoValidatorCommission(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeNoValidatorCommission), desc)
}

//func ErrWithdrawAddressInfoMismatch(codespace sdk.CodespaceType, expectedAddr, gotAddr sdk.AccAddress) sdk.Error {
//	return sdk.NewError(codespace, CodeErrWithdrawAddressInfoMismatch, fmt.Sprintf("account address mismatch, expected %s, got %s", expectedAddr.String(), gotAddr.String()))
//}
//
//func ErrHandleTxFailed(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeErrErrHandleTxFailed, err.Error())
//}
//
//func ErrInvalidSignature(codeSpace sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codeSpace, CodeErrInvalidSignature, msg)
//}