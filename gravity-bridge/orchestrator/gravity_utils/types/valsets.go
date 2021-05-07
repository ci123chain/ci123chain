package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type ValSet struct {
	Nonce uint64
	Members []ValSetMember
}

type ValSetMember struct {
	Power uint64
	EthAddress common.Address
}

