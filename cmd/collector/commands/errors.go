package commands

import "github.com/pkg/errors"

var (
	errMultipleAddFlags = errors.New("expected either --file/-f OR --url/u, found multiple")
)
