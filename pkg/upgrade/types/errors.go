package types

import sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"

var (
	ErrPermission = sdkerrors.Register(DefaultCodespace, 1707, "you have no permission to upgrade")
)
