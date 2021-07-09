package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCodec1(t *testing.T)  {
	acknowledgement := NewResultAcknowledgement([]byte{byte(1)})
	bz := acknowledgement.GetBytes()
	var ack Acknowledgement
	err := ChannelCdc.UnmarshalJSON(bz,&ack)
	require.Nil(t, err)
}
