package types
//
//import (
//	"fmt"
//	"strings"
//
//	cmn "github.com/ci123chain/ci123chain/pkg/libs/common"
//
//	abci "github.com/tendermint/tendermint/abci/types"
//)
//
// CodeType - ABCI code identifier within codespace
type CodeType uint32

// CodespaceType - codespace identifier
type CodespaceType string

//// IsOK - is everything okay?
func (code CodeType) IsOK() bool {
	if code == CodeOK {
		return true
	}
	return false
}

//const (
//	CodeOK    CodeType = 200
//)
//
//// SDK error codes
const (
	// Base error codes
	CodeOK                CodeType = 0
	//CodeInternal          CodeType = 1
	CodeTxDecode          CodeType = 2
	//CodeInvalidSequence   CodeType = 3
	//CodeUnauthorized      CodeType = 4
	//CodeInsufficientFunds CodeType = 5
	CodeUnknownRequest    CodeType = 6
	//CodeInvalidAddress    CodeType = 7
	//CodeInvalidPubKey     CodeType = 8
	//CodeUnknownAddress    CodeType = 9
	//CodeInsufficientCoins CodeType = 10
	//CodeInvalidCoins      CodeType = 11
	//CodeOutOfGas          CodeType = 12
	//CodeMemoTooLarge      CodeType = 13
	//CodeInsufficientFee   CodeType = 14
	//CodeTooManySignatures CodeType = 15
	//CodeNoPermission      CodeType = 16
	//CodeInvalidSign       CodeType = 17
	//
	//CodeNoDelegation      CodeType = 20
	//CodeNoValidator       CodeType = 21
	//CodeNoRedelegation    CodeType = 22
	//
	//CodeBadSharesAmount  CodeType = 30
	//
	//CodeValidateBlock CodeType = 40
	//CodeValidateEtherTx CodeType = 41

	// CodespaceRoot is a codespace for error codes in this file only.
	// Notice that 0 is an "unset" codespace, which can be overridden with
	// Error.WithDefaultCodespace().
	//CodespaceUndefined CodespaceType = ""
	CodespaceRoot      CodespaceType = "default"
)
//
//func unknownCodeMsg(code CodeType) string {
//	return fmt.Sprintf("unknown code %d", code)
//}
//
//// NOTE: Don't stringer this, we'll put better messages in later.
//func CodeToDefaultMsg(code CodeType) string {
//	switch code {
//	case CodeInternal:
//		return "internal error"
//	case CodeTxDecode:
//		return "tx parse error"
//	case CodeInvalidSequence:
//		return "invalid sequence"
//	case CodeUnauthorized:
//		return "unauthorized"
//	case CodeInsufficientFunds:
//		return "insufficient funds"
//	case CodeUnknownRequest:
//		return "unknown request"
//	case CodeInvalidAddress:
//		return "invalid address"
//	case CodeInvalidPubKey:
//		return "invalid pubkey"
//	case CodeUnknownAddress:
//		return "unknown address"
//	case CodeInsufficientCoins:
//		return "insufficient coins"
//	case CodeInvalidCoins:
//		return "invalid coins"
//	case CodeOutOfGas:
//		return "out of gas"
//	case CodeMemoTooLarge:
//		return "memo too large"
//	case CodeInsufficientFee:
//		return "insufficient fee"
//	case CodeTooManySignatures:
//		return "maximum numer of signatures exceeded"
//	case CodeValidateBlock:
//		return "validateBlock failed"
//	default:
//		return unknownCodeMsg(code)
//	}
//}
//
////--------------------------------------------------------------------------------
//// All errors are created via constructors so as to enable us to hijack them
//// and inject stack traces if we really want to.
//
//// nolint
//func ErrInternal(msg string) Error {
//	return newErrorWithRootCodespace(CodeInternal, msg)
//}
//func ErrTxDecode(msg string) Error {
//	return newErrorWithRootCodespace(CodeTxDecode, msg)
//}
//func ErrInvalidSequence(msg string) Error {
//	return newErrorWithRootCodespace(CodeInvalidSequence, msg)
//}
//func ErrUnauthorized(msg string) Error {
//	return newErrorWithRootCodespace(CodeUnauthorized, msg)
//}
//func ErrInsufficientFunds(msg string) Error {
//	return newErrorWithRootCodespace(CodeInsufficientFunds, msg)
//}
//func ErrUnknownRequest(msg string) Error {
//	return newErrorWithRootCodespace(CodeUnknownRequest, msg)
//}
//func ErrInvalidAddress(msg string) Error {
//	return newErrorWithRootCodespace(CodeInvalidAddress, msg)
//}
//func ErrUnknownAddress(msg string) Error {
//	return newErrorWithRootCodespace(CodeUnknownAddress, msg)
//}
//func ErrInvalidPubKey(msg string) Error {
//	return newErrorWithRootCodespace(CodeInvalidPubKey, msg)
//}
//func ErrInvalidSign(msg string) Error {
//	return newErrorWithRootCodespace(CodeInvalidSign, msg)
//}
//func ErrInsufficientCoins(msg string) Error {
//	return newErrorWithRootCodespace(CodeInsufficientCoins, msg)
//}
//func ErrInvalidCoins(msg string) Error {
//	return newErrorWithRootCodespace(CodeInvalidCoins, msg)
//}
//func ErrOutOfGas(msg string) Error {
//	return newErrorWithRootCodespace(CodeOutOfGas, msg)
//}
//func ErrMemoTooLarge(msg string) Error {
//	return newErrorWithRootCodespace(CodeMemoTooLarge, msg)
//}
//func ErrInsufficientFee(msg string) Error {
//	return newErrorWithRootCodespace(CodeInsufficientFee, msg)
//}
//func ErrTooManySignatures(msg string) Error {
//	return newErrorWithRootCodespace(CodeTooManySignatures, msg)
//}
//
//func ErrNoPermission(msg string) Error {
//	return newErrorWithRootCodespace(CodeNoPermission, msg)
//}
//
//func ErrNoDelegation(msg string) Error {
//	return newErrorWithRootCodespace(CodeNoDelegation, msg)
//}
//
//func ErrNoValidatorFound(msg string) Error {
//	return newErrorWithRootCodespace(CodeNoValidator, msg)
//}
//
//func ErrNoRedelegation(msg string) Error {
//	return newErrorWithRootCodespace(CodeNoRedelegation, msg)
//}
//
//func ErrBadSharesAmount(msg string) Error {
//	return newErrorWithRootCodespace(CodeBadSharesAmount, msg)
//}
//
//func ErrValidateBlock(msg string) Error {
//	return newErrorWithRootCodespace(CodeValidateBlock, msg)
//}
//
//func ErrValidateEtherTx(msg string) Error {
//	return newErrorWithRootCodespace(CodeValidateEtherTx, msg)
//}
//
////----------------------------------------
//// Error & sdkError
//
//type cmnError = cmn.Error
//
//// sdk Error type
//type Error interface {
//	// Implements cmn.Error
//	// Error() string
//	// Stacktrace() cmn.Error
//	// Trace(offset int, format string, args ...interface{}) cmn.Error
//	// Data() interface{}
//	cmnError
//
//	// convenience
//	TraceSDK(format string, args ...interface{}) Error
//
//	// set codespace
//	WithDefaultCodespace(CodespaceType) Error
//
//	Code() CodeType
//	Codespace() CodespaceType
//	ABCILog() string
//	Result() Result
//	QueryResult() abci.ResponseQuery
//}
//
//// NewError - create an error.
//func NewError(codespace CodespaceType, code CodeType, format string, args ...interface{}) Error {
//	return newError(codespace, code, format, args...)
//}
//
//func newErrorWithRootCodespace(code CodeType, format string, args ...interface{}) *sdkError {
//	return newError(CodespaceRoot, code, format, args...)
//}
//
//func newError(codespace CodespaceType, code CodeType, format string, args ...interface{}) *sdkError {
//	if format == "" {
//		format = CodeToDefaultMsg(code)
//	}
//	return &sdkError{
//		codespace: codespace,
//		code:      code,
//		cmnError:  cmn.NewError(format, args...),
//	}
//}
//
//type sdkError struct {
//	codespace CodespaceType
//	code      CodeType
//	cmnError
//}
//
//// Implements Error.
//func (err *sdkError) WithDefaultCodespace(cs CodespaceType) Error {
//	codespace := err.codespace
//	if codespace == CodespaceUndefined {
//		codespace = cs
//	}
//	return &sdkError{
//		codespace: cs,
//		code:      err.code,
//		cmnError:  err.cmnError,
//	}
//}
//
//// Implements ABCIError.
//// nolint: errcheck
//func (err *sdkError) TraceSDK(format string, args ...interface{}) Error {
//	_ = err.Trace(1, format, args...)
//	return err
//}
//
//// Implements ABCIError.
//func (err *sdkError) Error() string {
//	return fmt.Sprintf(`ERROR:
//Codespace: %s
//Code: %d
//Message: %#v
//`, err.codespace, err.code, err.cmnError.Error())
//}
//
//// Implements Error.
//func (err *sdkError) Codespace() CodespaceType {
//	return err.codespace
//}
//
//// Implements Error.
//func (err *sdkError) Code() CodeType {
//	return err.code
//}
//
//// Implements ABCIError.
//func (err *sdkError) ABCILog() string {
//	//cdc := codec.New()
//	//errMsg := err.cmnError.Error()
//	//jsonErr := humanReadableError{
//	//	Codespace: err.codespace,
//	//	Code:      err.code,
//	//	Message:   errMsg,
//	//}
//	//bz, er := cdc.MarshalJSON(jsonErr)
//	//if er != nil {
//	//	panic(er)
//	//}
//	//stringifiedJSON := string(bz)
//	//return stringifiedJSON
//
//	return err.cmnError.Error()
//}
//
//func (err *sdkError) Result() Result {
//	return Result{
//		Code:      err.Code(),
//		Codespace: err.Codespace(),
//		Log:       err.ABCILog(),
//	}
//}
//
//// QueryResult allows us to return sdk.Error.QueryResult() in query responses
//func (err *sdkError) QueryResult() abci.ResponseQuery {
//	return abci.ResponseQuery{
//		Code:      uint32(err.Code()),
//		Codespace: string(err.Codespace()),
//		Log:       err.ABCILog(),
//	}
//}
//
////----------------------------------------
//// REST error utilities
//
//// appends a message to the head of the given error
//func AppendMsgToErr(msg string, err string) string {
//	msgIdx := strings.Index(err, "message\":\"")
//	if msgIdx != -1 {
//		errMsg := err[msgIdx+len("message\":\"") : len(err)-2]
//		errMsg = fmt.Sprintf("%s; %s", msg, errMsg)
//		return fmt.Sprintf("%s%s%s",
//			err[:msgIdx+len("message\":\"")],
//			errMsg,
//			err[len(err)-2:],
//		)
//	}
//	return fmt.Sprintf("%s; %s", msg, err)
//}
//
//// returns the index of the message in the ABCI Log
//func mustGetMsgIndex(abciLog string) int {
//	msgIdx := strings.Index(abciLog, "message\":\"")
//	if msgIdx == -1 {
//		panic(fmt.Sprintf("invalid error format: %s", abciLog))
//	}
//	return msgIdx + len("message\":\"")
//}
//
//// parses the error into an object-like struct for exporting
//type humanReadableError struct {
//	Codespace CodespaceType `json:"codespace"`
//	Code      CodeType      `json:"code"`
//	Message   string        `json:"message"`
//}
