package types

import "github.com/pkg/errors"

var (
	ErrInvalidRoute = errors.New("port route invalid")
	ErrInvalidPort = errors.New("port id invalid")
)
