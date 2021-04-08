package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace 				sdk.CodespaceType = "auth"
	CodeTxValidateBasicError       	CodeType = 501
)

func ErrTxValidateBasic(codespace sdk.CodespaceType, err error) error {
	return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("validate basic eror: %v", err.Error()))
}
