package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
// SubModuleName is the error codespace
const SubModuleName string = "commitment"

// IBC connection sentinel errors
var (
	ErrInvalidProof       = sdkerrors.Register(SubModuleName, 2072, "invalid proof")
	ErrInvalidPrefix      = sdkerrors.Register(SubModuleName, 2073, "invalid prefix")
	ErrInvalidMerkleProof = sdkerrors.Register(SubModuleName, 2074, "invalid merkle proof")
)
