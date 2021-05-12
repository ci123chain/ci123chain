package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace  						sdk.CodespaceType = "mint"

	CodeInvalidEndPoint                    CodeType = 1801
)
//
//func ErrFailedMarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
//	return sdk.NewError(codespace, CodeBadMarshal, "Marshal Error: %s", detailStr)
//}

func ErrInvalidEndPoint(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidEndPoint), desc)
}