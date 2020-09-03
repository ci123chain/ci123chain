package types

import (
	"fmt"
	"github.com/pkg/errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace  sdk.CodespaceType = "staking"
	StakingCodespace = "staking"

	CodeValidatorExisted CodeType = 66

	CodeDescriptionOutOfLength CodeType = 68
	CodeSetCommissionFailed CodeType = 69
	CodeSetValidatorFailed CodeType = 70
	CodeDelegateFailed CodeType = 71
	CodeNoExpectedValidator CodeType = 72
	CodeUnexpectedDeom CodeType = 73
	CodeRedelegateFailed CodeType = 74
	CodeGotTimeFailed CodeType = 75
	CodeUndelegateFailed CodeType = 76
	CodeValidateUnbondAmountFailed CodeType = 77

	CodeCheckParamsError	CodeType = 79
	CodeInvalidAddress      CodeType = 80
	CodeEmptyPublicKey      CodeType = 81
)

var (
	ErrCommissionUpdateTime = Register(StakingCodespace, 17, "new rate cannot be changed more than once within 24 hours")
	ErrInvalidRequest = Register(StakingCodespace, 18, "invalid request")
	ErrCommissionNegative = Register(StakingCodespace, 19, "commission must be positive" )
	ErrCommissionHuge = Register(StakingCodespace, 20, "commission cannot be more than 100%")
	ErrCommissionGTMaxRate = Register(StakingCodespace, 21, "commission change rate cannot be more than the max rate")
	ErrCommissionChangeRateNegative = Register(StakingCodespace, 22, "commission change rate must be positive")
	ErrCommissionChangeRateGTMaxRate = Register(StakingCodespace, 23, "commission change rate cannot be more than the max rate")
	ErrDelegatorShareExRateInvalid = Register(StakingCodespace, 24, "cannot delegate to validators with invalid (zero) ex-rate")
	ErrInsufficientShares = Register(StakingCodespace, 25, "insufficient delegation shares")
	ErrUnknowTokenSource = Register(StakingCodespace, 26, "unknown token source bond status")
	ErrInvalidValidatorStatus = Register(StakingCodespace, 27, "invalid validator status")
	ErrBondedTokendFailed = Register(StakingCodespace, 28, "notBondedTokensToBonded failed")
	ErrBondedTokensToNoBondedFailed = Register(StakingCodespace, 29, "BondedTokensToBonded failed")
	ErrNoDelegation = Register(StakingCodespace, 30, "no delegation")
	ErrBadSharesAmount = Register(StakingCodespace, 31, "invalid shares amount")
	ErrSelfRedelegation = Register(StakingCodespace, 32, "cannot redelegate to the same validator")
	ErrNoValidatorFound = Register(StakingCodespace, 33, "no validator found")
	ErrNotEnoughDelegationShares = Register(StakingCodespace, 34, "not enough delegation shares")
	ErrNoDelegatorForAddress = Register(StakingCodespace, 35, "delegator does not contain delegation")
	ErrMaxUnbondingDelegationEntries = Register(StakingCodespace, 36, "too many unbonding delegation entries for (delegator, validator) tuple")
	ErrTinyRedelegationAmount = Register(StakingCodespace, 37, "too few tokens to redelegate (truncates to zero tokens)")
	ErrMaxRedelegationEntries = Register(StakingCodespace, 38, "too many redelegation entries for (delegator, src-validator, dst-validator) tuple")
	ErrTransitiveRedelegation = Register(StakingCodespace, 39, "redelegation to this validator already in progress; first redelegation to this validator must complete before next redelegation")
	ErrBadRedelegationDst = Register(StakingCodespace, 40, "redelegation destination validator not found")
)

func Wrapf(err error, format string, args ...interface{}) error {
	desc := fmt.Sprintf(format, args...)
	return Wrap(err, desc)
}

func Wrap(err error, description string) error {
	if err == nil {
		return nil
	}
	// If this error does not carry the stacktrace information yet, attach
	// one. This should be done only once per error at the lowest frame
	// possible (most inner wrap).
	if stackTrace(err) == nil {
		err = errors.WithStack(err)
	}

	return &wrappedError{
		parent: err,
		msg:    description,
	}
}

type wrappedError struct {
	// This error layer description.
	msg string
	// The underlying error that triggered this one.
	parent error
}

func (e *wrappedError) Error() string {
	return fmt.Sprintf("%s: %s", e.parent.Error(), e.msg)
}

func stackTrace(err error) errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	for {
		if st, ok := err.(stackTracer); ok {
			return st.StackTrace()
		}

		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			return nil
		}
	}
}

// causer is an interface implemented by an error that supports wrapping. Use
// it to test if an error wraps another error instance.
type causer interface {
	Cause() error
}

type unpacker interface {
	Unpack() []error
}


type Error struct {
	codespace string
	code      uint32
	desc      string
}

func (e Error) Error() string {
	return e.desc
}

func New(codespace string, code uint32, desc string) *Error {
	return &Error{codespace: codespace, code: code, desc: desc}
}


func Register(codespace string, code uint32, description string) *Error {

	// TODO - uniqueness is (codespace, code) combo
	if e := getUsed(codespace, code); e != nil {
		panic(fmt.Sprintf("error with code %d is already registered: %q", code, e.desc))
	}

	err := New(codespace, code, description)
	setUsed(err)

	return err
}

// usedCodes is keeping track of used codes to ensure their uniqueness. No two
// error instances should share the same (codespace, code) tuple.
var usedCodes = map[string]*Error{}

func errorID(codespace string, code uint32) string {
	return fmt.Sprintf("%s:%d", codespace, code)
}

func getUsed(codespace string, code uint32) *Error {
	return usedCodes[errorID(codespace, code)]
}

func setUsed(err *Error) {
	usedCodes[errorID(err.codespace, err.code)] = err
}

func ErrValidatorExisted(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeValidatorExisted, "Validator existed: %s", err.Error())
}

func ErrDescriptionOutOfLength(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeDescriptionOutOfLength, "Description out of length: %s", err.Error())
}

func ErrSetCommissionFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeSetCommissionFailed, "Set commission failed: %s", err.Error())
}
func ErrSetValidatorFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeSetValidatorFailed, "Set validator failed: %s", err.Error())
}

func ErrDelegateFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeDelegateFailed, "Delegate failed: %s", err.Error())
}

func ErrNoExpectedValidator(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeNoExpectedValidator, "No expected validator: %s", err.Error())
}

func ErrBondedDenomDiff(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeUnexpectedDeom, "unexpected denom: %s", err.Error())
}

func ErrRedelegationFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeRedelegateFailed, "Redelegatie failed: %s", err.Error())
}

func ErrGotTimeFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGotTimeFailed, "got time failed: %s", err.Error())
}

func ErrUndelegateFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeUndelegateFailed, "Undelegate failed: %s", err.Error())
}

func ErrValidateUnBondAmountFailed(codespace sdk.CodespaceType, err error) sdk.Error {

	return sdk.NewError(codespace, CodeValidateUnbondAmountFailed, "Validate unbond amount failed: %s", err.Error())
}

func ErrCheckParams(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid: %s", str)
}

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

func ErrEmptyPublicKey(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyPublicKey, msg)
}