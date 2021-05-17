package net

import (
	"fmt"
	clientcontext "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/util"
)

// PublicNetAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicNetAPI struct {
	networkVersion uint64
}

// NewAPI creates an instance of the public Net Web3 API.
func NewAPI(clientCtx clientcontext.Context) *PublicNetAPI {
	return &PublicNetAPI{
		networkVersion: uint64(util.CHAINID),
	}
}

// Version returns the current ethereum protocol version.
func (api *PublicNetAPI) Version() string {
	return fmt.Sprintf("%d", api.networkVersion)
}

// Listening returns an indication if the node is listening for network connections.
func (api *PublicNetAPI) Listening() bool {
	return true // always listening
}