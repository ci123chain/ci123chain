package types
//
import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
//
type CodeType = sdk.CodeType
const (
	DefaultCodespace 				sdk.CodespaceType = "client"
	//CodeNewClientCtxError       	CodeType = 601
	//CodeGetInputAddrError       	CodeType = 602
	//CodeParseAddrError       		CodeType = 603
	//CodeNoAddrError       			CodeType = 604
	//CodeGetPassPhraseError			CodeType = 605
	//CodeGetSignDataError			CodeType = 606
	//CodeBroadcastError				CodeType = 607
	//CodeGetCheckPasswordError		CodeType = 608
	//CodeGetPasswordError			CodeType = 609
	//CodePhrasesNotMatchError		CodeType = 610
	//CodeNodeError					CodeType = 611
	//CodeParseParamsError       		CodeType = 612
	//CodeGenValidatorErr             CodeType = 613
	//CodeGenAccountErr               CodeType = 614
)
//
//func ErrNewClientCtx(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeNewClientCtxError,"New client context Error: %s", err.Error())
//}
//
//func ErrGetInputAddrCtx(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGetInputAddrError,"Get input address Error: %s", err.Error())
//}
//
//func ErrParseAddr(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeParseAddrError,"Parse address Error: %s", err.Error())
//}
//
//func ErrNoAddr(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeNoAddrError,"No address Error: %s", err.Error())
//}
//
//func ErrGetPassPhrase(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGetPassPhraseError,"Get pass phrase Error: %s", err.Error())
//}
//
//func ErrGetSignData(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGetSignDataError,"Get sign data from tx Error: %s", err.Error())
//}
//
//func ErrBroadcast(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeBroadcastError,"Broadcast Error: %s", err.Error())
//}
//
//func ErrGetCheckPassword(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGetCheckPasswordError,"Get check password Error: %s", err.Error())
//}
//
//func ErrGetPassword(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGetPasswordError,"Get password Error: %s", err.Error())
//}
//
//func ErrPhrasesNotMatch(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodePhrasesNotMatchError,"Phrases not match Error: %s", err.Error())
//}
//
func ErrNode(codespace sdk.CodespaceType, err error) error {
	return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("get node failed: %v", err))
}
//
//func ErrGenValidatorKey(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGenValidatorErr,"Gen Validator Key Error: %s", err.Error())
//}
//
//
//func ErrParseParam(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeParseParamsError,"Parse params Error: %s", err.Error())
//}
//
//func ErrGenAccount(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGenAccountErr,"Gen account: %s", err.Error())
//}