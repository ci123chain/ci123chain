package types

import "github.com/pkg/errors"

var (
	ErrProcessedTimeNotFound = errors.New("processed time not found !")
	ErrDelayPeriodNotPassed = errors.New("delay period not passed")
)