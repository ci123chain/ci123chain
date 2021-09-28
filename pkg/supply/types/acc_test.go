package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"testing"
)

func TestAcc(t *testing.T) {
	a := types.ToAccAddress(crypto.AddressHash([]byte("preStaking")))
	fmt.Println(a.String())
}
