package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
)

type Erc20Token struct {
	Amount sdk.Int
	TokenContractAddress common.Address
}


