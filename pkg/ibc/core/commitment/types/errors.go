package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
// SubModuleName is the error codespace
const SubModuleName string = "commitment"

// IBC connection sentinel errors
var (
	ErrInvalidProof       = sdkerrors.Register(SubModuleName, 2, "invalid proof")
	ErrInvalidPrefix      = sdkerrors.Register(SubModuleName, 3, "invalid prefix")
	ErrInvalidMerkleProof = sdkerrors.Register(SubModuleName, 4, "invalid merkle proof")
)
