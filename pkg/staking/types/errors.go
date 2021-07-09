package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace  sdk.CodespaceType = "staking"
	StakingCodespace = "staking"

	//CodeValidatorExisted CodeType = 66
	//
	//CodeDescriptionOutOfLength CodeType = 68
	//CodeSetCommissionFailed CodeType = 69
	//CodeSetValidatorFailed CodeType = 70
	//CodeDelegateFailed CodeType = 71
	//CodeNoExpectedValidator CodeType = 72
	//CodeUnexpectedDeom CodeType = 73
	//CodeRedelegateFailed CodeType = 74
	//CodeGotTimeFailed CodeType = 75
	//CodeUndelegateFailed CodeType = 76
	//CodeValidateUnbondAmountFailed CodeType = 77
	//
	//CodeCheckParamsError	CodeType = 79
	//CodeInvalidAddress      CodeType = 80
	//CodeEmptyPublicKey      CodeType = 81

	CodeInvalidValidator CodeType = 1601
	CodeInvalidTxType    CodeType = 1602
	CodeInvalidPublicKey CodeType = 1603
	CodeInvalidParam     CodeType = 1604

	CodeInternalOperationFailed CodeType = 1605
	Code

	CodeModuleAccountNotExisted CodeType = 1606
)

var (
	ErrCommissionUpdateTime = sdkerrors.Register(StakingCodespace, 1607, "new rate cannot be changed more than once within 24 hours")
	ErrInvalidRequest = sdkerrors.Register(StakingCodespace, 1608, "invalid request")
	ErrCommissionNegative = sdkerrors.Register(StakingCodespace, 1609, "commission must be positive" )
	ErrCommissionHuge = sdkerrors.Register(StakingCodespace, 1610, "commission cannot be more than 100%")
	ErrCommissionGTMaxRate = sdkerrors.Register(StakingCodespace, 1611, "commission change rate cannot be more than the max rate")
	ErrCommissionChangeRateNegative = sdkerrors.Register(StakingCodespace, 1612, "commission change rate must be positive")
	ErrCommissionChangeRateGTMaxRate = sdkerrors.Register(StakingCodespace, 1613, "commission change rate cannot be more than the max rate")
	ErrDelegatorShareExRateInvalid = sdkerrors.Register(StakingCodespace, 1614, "cannot delegate to validators with invalid (zero) ex-rate")
	ErrInsufficientShares = sdkerrors.Register(StakingCodespace, 1615, "insufficient delegation shares")
	ErrUnknowTokenSource = sdkerrors.Register(StakingCodespace, 1616, "unknown token source bond status")
	ErrInvalidValidatorStatus = sdkerrors.Register(StakingCodespace, 1617, "invalid validator status")
	ErrBondedTokendFailed = sdkerrors.Register(StakingCodespace, 1618, "notBondedTokensToBonded failed")
	ErrBondedTokensToNoBondedFailed = sdkerrors.Register(StakingCodespace, 1619, "BondedTokensToBonded failed")
	ErrNoDelegation = sdkerrors.Register(StakingCodespace, 1620, "no delegation")
	ErrBadSharesAmount = sdkerrors.Register(StakingCodespace, 1621, "invalid shares amount")
	ErrSelfRedelegation = sdkerrors.Register(StakingCodespace, 1622, "cannot redelegate to the same validator")
	ErrNoValidatorFound = sdkerrors.Register(StakingCodespace, 1623, "no validator found")
	ErrNotEnoughDelegationShares = sdkerrors.Register(StakingCodespace, 1624, "not enough delegation shares")
	ErrNoDelegatorForAddress = sdkerrors.Register(StakingCodespace, 1625, "delegator does not contain delegation")
	ErrMaxUnbondingDelegationEntries = sdkerrors.Register(StakingCodespace, 1626, "too many unbonding delegation entries for (delegator, validator) tuple")
	ErrTinyRedelegationAmount = sdkerrors.Register(StakingCodespace, 1627, "too few tokens to redelegate (truncates to zero tokens)")
	ErrMaxRedelegationEntries = sdkerrors.Register(StakingCodespace, 1628, "too many redelegation entries for (delegator, src-validator, dst-validator) tuple")
	ErrTransitiveRedelegation = sdkerrors.Register(StakingCodespace, 1629, "redelegation to this validator already in progress; first redelegation to this validator must complete before next redelegation")
	ErrBadRedelegationDst = sdkerrors.Register(StakingCodespace, 1630, "redelegation destination validator not found")
	ErrCdcMarshal = sdkerrors.Register(StakingCodespace, 1631, "codec marshal failed")
	ErrCdcUnmarshal = sdkerrors.Register(StakingCodespace, 1632, "codec unmarshal failed")

	ErrNoExpectedValidator = sdkerrors.Register(StakingCodespace, uint32(CodeInvalidValidator),  "no expected validator found")
	ErrInvalidTxType = sdkerrors.Register(StakingCodespace, uint32(CodeInvalidTxType), "invalid tx type")
	ErrInvalidPublicKey = sdkerrors.Register(StakingCodespace, uint32(CodeInvalidPublicKey), "invalid publickey")
	ErrPubkeyHasBonded = sdkerrors.Register(StakingCodespace, uint32(1634), "pubkey has bonded with one account")
	ErrInvalidMoniker = sdkerrors.Register(StakingCodespace, uint32(CodeInvalidParam), "invalid moniker")
	ErrInvalidIdentity = sdkerrors.Register(StakingCodespace, uint32(1635), "invalid identity")
	ErrInvalidWebsite = sdkerrors.Register(StakingCodespace, uint32(1636), "invalid website")
	ErrInvalidSecurityContact = sdkerrors.Register(StakingCodespace, uint32(1637), "invalid security contact")
	ErrInvalidDetails = sdkerrors.Register(StakingCodespace, uint32(1638), "invalid details")
	ErrSetValidatorFailed = sdkerrors.Register(StakingCodespace, uint32(CodeInternalOperationFailed), "set validator failed")
	ErrInvalidParam = sdkerrors.Register(StakingCodespace, uint32(1639), "invalid params")
)

//func Wrapf(err error, format string, args ...interface{}) error {
//	desc := fmt.Sprintf(format, args...)
//	return Wrap(err, desc)
//}
//
//func Wrap(err error, description string) error {
//	if err == nil {
//		return nil
//	}
//	// If this error does not carry the stacktrace information yet, attach
//	// one. This should be done only once per error at the lowest frame
//	// possible (most inner wrap).
//	if stackTrace(err) == nil {
//		err = errors.WithStack(err)
//	}
//
//	return &wrappedError{
//		parent: err,
//		msg:    description,
//	}
//}
//
//type wrappedError struct {
//	// This error layer description.
//	msg string
//	// The underlying error that triggered this one.
//	parent error
//}
//
//func (e *wrappedError) Error() string {
//	return fmt.Sprintf("%s: %s", e.parent.Error(), e.msg)
//}
//
//func stackTrace(err error) errors.StackTrace {
//	type stackTracer interface {
//		StackTrace() errors.StackTrace
//	}
//
//	for {
//		if st, ok := err.(stackTracer); ok {
//			return st.StackTrace()
//		}
//
//		if c, ok := err.(causer); ok {
//			err = c.Cause()
//		} else {
//			return nil
//		}
//	}
//}
//
//// causer is an interface implemented by an error that supports wrapping. Use
//// it to test if an error wraps another error instance.
//type causer interface {
//	Cause() error
//}
//
//type unpacker interface {
//	Unpack() []error
//}
//
//
//type Error struct {
//	codespace string
//	code      uint32
//	desc      string
//}
//
//func (e Error) Error() string {
//	return e.desc
//}
//
//func New(codespace string, code uint32, desc string) *Error {
//	return &Error{codespace: codespace, code: code, desc: desc}
//}
//
//
//func Register(codespace string, code uint32, description string) *Error {
//
//	// TODO - uniqueness is (codespace, code) combo
//	if e := getUsed(codespace, code); e != nil {
//		panic(fmt.Sprintf("error with code %d is already registered: %q", code, e.desc))
//	}
//
//	err := New(codespace, code, description)
//	setUsed(err)
//
//	return err
//}
//
//// usedCodes is keeping track of used codes to ensure their uniqueness. No two
//// error instances should share the same (codespace, code) tuple.
//var usedCodes = map[string]*Error{}
//
//func errorID(codespace string, code uint32) string {
//	return fmt.Sprintf("%s:%d", codespace, code)
//}
//
//func getUsed(codespace string, code uint32) *Error {
//	return usedCodes[errorID(codespace, code)]
//}
//
//func setUsed(err *Error) {
//	usedCodes[errorID(err.codespace, err.code)] = err
//}

//func ErrValidatorExisted(codespace sdk.CodespaceType, _ error) sdk.Error {
//	return sdk.NewError(codespace, CodeValidatorExisted, "Validator existed already", "")
//}

//func ErrDescriptionOutOfLength(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeDescriptionOutOfLength, "Description out of length: %s", err.Error())
//}

//func ErrSetCommissionFailed(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeSetCommissionFailed, "Set commission failed: %s", err.Error())
//}
//func ErrSetValidatorFailed(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeSetValidatorFailed, "Set validator failed: %s", err.Error())
//}

//func ErrDelegateFailed(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeDelegateFailed, "Delegate fail: %s", err.Error())
//}

//func ErrNoExpectedValidator(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidValidator),  desc)
//}
//
//func ErrBondedDenomDiff(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeUnexpectedDeom, "unexpected denom: %s", err.Error())
//}

//func ErrRedelegationFailed(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeRedelegateFailed, "Redelegatie failed: %s", err.Error())
//}

//func ErrGotTimeFailed(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGotTimeFailed, "got time failed: %s", err.Error())
//}

//func ErrUndelegateFailed(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeUndelegateFailed, "Undelegate failed: %s", err.Error())
//}

//func ErrValidateUnBondAmountFailed(codespace sdk.CodespaceType, err error) sdk.Error {
//
//	return sdk.NewError(codespace, CodeValidateUnbondAmountFailed, "Validate unbond amount failed: %s", err.Error())
//}

//func ErrCheckParams(codespace sdk.CodespaceType, str string) sdk.Error {
//	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid: %s", str)
//}

//func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespace, CodeInvalidAddress, msg)
//}

//func ErrEmptyPublicKey(codespace sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespace, CodeEmptyPublicKey, msg)
//}

//func ErrInvalidTxType(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidTxType), desc)
//}

//func ErrInvalidPublicKey(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidPublicKey), desc)
//}

//func ErrPubkeyHasBonded(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidPublicKey), desc)
//}

//func ErrInvalidMoniker(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidParam), desc)
//}

//func ErrInvalidIdentity(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidParam), desc)
//}

//func ErrInvalidWebsite(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidParam), desc)
//}

//func ErrInvalidSecurityContact(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidParam), desc)
//}

//func ErrInvalidDetails(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidParam), desc)
//}


//func ErrSetValidatorFailed(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInternalOperationFailed), desc)
//}

//func ErrInvalidParam(desc string) error {
//	return sdkerrors.Register(StakingCodespace, uint32(CodeInvalidParam), desc)
//}