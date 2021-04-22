package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.HexToAddress("abcd")
	msg := NewMsgUnjail(sdk.AccAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"ci123chain/MsgUnjail","value":{"address":"cosmosvaloper1v93xxeqhg9nn6"}}`,
		string(bytes),
	)
}
