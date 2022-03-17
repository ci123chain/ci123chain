package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

var (
	ErrInternal                = sdkerrors.Register(ModuleName, 1, "internal")
	ErrDuplicate               = sdkerrors.Register(ModuleName, 2, "duplicate")
	ErrInvalid                 = sdkerrors.Register(ModuleName, 3, "invalid")
	ErrTimeout                 = sdkerrors.Register(ModuleName, 4, "timeout")
	ErrUnknown                 = sdkerrors.Register(ModuleName, 5, "unknown")
	ErrEmpty                   = sdkerrors.Register(ModuleName, 6, "empty")
	ErrOutdated                = sdkerrors.Register(ModuleName, 7, "outdated")
	ErrUnsupported             = sdkerrors.Register(ModuleName, 8, "unsupported")
	ErrNonContiguousEventNonce = sdkerrors.Register(ModuleName, 9, "non contiguous event nonce")
	ErrNoContractMetaData 	   = sdkerrors.Register(ModuleName, 10, "contract metadata not found")
	ErrMappedContractNotFound  = sdkerrors.Register(ModuleName, 11, "mapped contract not found")
	ErrBatchIDNil  			   = sdkerrors.Register(ModuleName, 12, "batch id can not be nil")

	ErrDenomDecimal  		   = sdkerrors.Register(ModuleName, 13, "denom decimal invalid")
	ErrDenomName  			   = sdkerrors.Register(ModuleName, 14, "denom name invalid")
	ErrDenomSymbol			   = sdkerrors.Register(ModuleName, 15, "denom symbol invalid")

	ErrQueryDenom 			   = sdkerrors.Register(ModuleName, 16, "denom query failed")
	ErrQueryDenomMismatch 	   = sdkerrors.Register(ModuleName, 17, "denom query mismatch")

	ErrQueryERC20 			   = sdkerrors.Register(ModuleName, 18, "erc20 query failed")
	ErrNoTxToRelay 			   = sdkerrors.Register(ModuleName, 19, "no such tx to request batch")

)
