package types

import (
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
)

func (m Attestation) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(m.Claim, new(EthereumClaim))
}