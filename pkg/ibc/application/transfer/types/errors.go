package types

import "github.com/pkg/errors"

var (
	ErrInvalidAmount = errors.New("invalid amount !")
	ErrInvalidAddress = errors.New("invalid address !")
	ErrSendDisabled = errors.New("send disable !")
	ErrReceiveDisabled = errors.New("receive disable !")
	ErrInvalidDenomForTransfer = errors.New("invalid denom for transfer !")
	ErrTraceNotFound = errors.New("err trace not found !")
	ErrUnknownRequest = errors.New("unknow request !")
	ErrInvalidVersion = errors.New("countparty version invalid !")
	ErrMaxTransferChannels = errors.New("channel sequence bigger than max")

)
