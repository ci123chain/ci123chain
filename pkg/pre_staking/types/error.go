package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type CodeType = sdk.CodeType



var (
	ErrAccountBalanceNotEnough = sdkerrors.Register(DefaultCodespace, 1800, "balance of account not enough to pay")
	ErrInvalidAmount = sdkerrors.Register(DefaultCodespace, 1801, "invalid amount")
	ErrInvalidDelegatorAddress = sdkerrors.Register(DefaultCodespace, 1802, "invalid delegator address")
	ErrInvalidValidatorAddress = sdkerrors.Register(DefaultCodespace, 1803, "invalid validator address")
	ErrFromNotEqualDelegator = sdkerrors.Register(DefaultCodespace, 1804, "from address not equal to delegator address")
	ErrNoExpectedValidator = sdkerrors.Register(DefaultCodespace, 1805, "no expected validator found")
	ErrInvalidDenom = sdkerrors.Register(DefaultCodespace, 1806, "invalid denom")
	ErrDelegatorShareExRateInvalid = sdkerrors.Register(DefaultCodespace, 1807, "invalid delegator share exchange rate")
)