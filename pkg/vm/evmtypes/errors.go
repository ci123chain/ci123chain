package evmtypes
//
import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
//
type CodeType = sdk.CodeType
const (
	DefaultCodespace sdk.CodespaceType = "evm"

	CodeCheckParamsError	CodeType = 1701

	CodeComputeGas          CodeType = 1702
	CodeInvalidState        CodeType = 1703
	CodeMissingPassword     CodeType = 1704
	CodeInvalidPassword     CodeType = 1705
	CodeTimeOut             CodeType = 1706
	CodeInvalidPrivateKey   CodeType = 1707
	CodeUnlockTimeTooLarge  CodeType = 1708
	CodeInvalidClientStatus CodeType = 1709
	CodeGetNodeFailed       CodeType = 1710
	CodeInvalidNodeStatus  CodeType = 1711
	CodeGetABCIInfoFailed   CodeType = 1712
	CodeGetBlockNumber      CodeType = 1713
	CodeGetBlockFailed      CodeType = 1714
	CodeCdcMarshalFailed    CodeType = 1715
	CodeCdcUnmarshalFailed  CodeType = 1716
	CodeeQueryAccountFailed  CodeType = 1717
	CodeGetStorageAtFailed   CodeType = 1718
	CodeClientQueryTxFailed  CodeType = 1719
	CodeAccountExisted       CodeType = 1720
	CodeImportRawKeyFailed   CodeType = 1721
)

var (
	ErrComputeGas = sdkerrors.Register(string(DefaultCodespace), uint32(CodeComputeGas), "invalid intrinsic gas for transaction")
	ErrInvalidState = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidState), "invalid vm state")
	ErrInvalidParams = sdkerrors.Register(string(DefaultCodespace), uint32(CodeMissingPassword), "invalid params")
	ErrMissingPassword = sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), "missing password")
	ErrInvalidPassword = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidPassword), "invalid password")
	ErrTimeOut = sdkerrors.Register(string(DefaultCodespace), uint32(CodeTimeOut), "timeout")
	ErrInvalidPrivateKey= sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidPrivateKey), "invalid private_key")
	ErrUnlockTimeTooLarger = sdkerrors.Register(string(DefaultCodespace), uint32(CodeUnlockTimeTooLarge), "unlock time too large")
	ErrInvalidClientStatus = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidClientStatus), "unlock time too large")
	ErrGetNodeFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetNodeFailed), "get node failed")
	ErrInvalidNodeStatus = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidNodeStatus), "invalid node status")
	ErrGetABCIInfoFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetABCIInfoFailed), "get abci info failed")
	ErrGetBlockNumber = sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetBlockNumber), "get block number failed")
	ErrGetBlockFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetBlockFailed), "get block failed")
	ErrCdcMarshalFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcMarshalFailed), "cdc marshal failed")
	ErrCdcUnmarshalFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcUnmarshalFailed), "cdc unmarshal failed")
	ErrQueryAccountsFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeeQueryAccountFailed), "query account failed")
	ErrGetStorageAtFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetStorageAtFailed), "get storage failed")
	ErrClientQueryTxFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeClientQueryTxFailed), "client query tx failed")
	ErrAccountExisted = sdkerrors.Register(string(DefaultCodespace), uint32(CodeAccountExisted), "account already exists")
	ErrImportRawKeyFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeImportRawKeyFailed), "import raw key failed")
)
//
//const (
//	CodeCheckParamsError	CodeType = 60
//	CodeInvalidMsgError     CodeType = 61
//	CodeHandleMsgFailedError  CodeType = 62
//	CodeSetSequenceFailedError CodeType = 63
//	CodeInvalidAddress        CodeType  = 64
//	CodeComputeGasFailedError  CodeType  = 65
//	CodeErrInvalidState			CodeType  = 66
//	CodeErrChainConfigNotFound CodeType = 67
//	CodeErrTransitionDb		CodeType = 68
//)
//
//func ErrCheckParams(codespace sdk.CodespaceType, keyname string) sdk.Error {
//	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid: %s", keyname)
//}
//
//func ErrInvalidMsg(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeInvalidMsgError, "msg invalid: %s", err.Error())
//}
//
//func ErrSetNewAccountSequence(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeSetSequenceFailedError, "set sequence of account failed: %s", err.Error())
//}
//
//func ErrInvalidAddress(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeInvalidAddress, msg)
//}
//
//func ErrComputeGas(desc string) error {
//	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeComputeGas), desc)
//}
////
//func ErrInvalidState(desc string) error {
//	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidState), desc)
//}
//
//func ErrInvalidParams(desc string) error {
//	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), desc)
//}
//
//func ErrChainConfigNotFound(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeErrChainConfigNotFound, msg)
//}
//
//func ErrTransitionDb(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeErrTransitionDb, msg)
//}
//
