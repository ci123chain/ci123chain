package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnection(t *testing.T)  {
	addr := "https://tm.weelinknode1c.gw002.oneitfarm.com:443"
	_, err := GetConnection(addr)
	assert.NoError(t, err)
}
