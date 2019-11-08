package sdk

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"testing"
)

var testPriv2 = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg8j+IqOSi1dVmexs5
8Vx+8knGdeibNvwMTyZ05QO32LmhRANCAAR3LEnYWIDz9orUF2v7wES2vjxArABy
neS3btroKyAsjEXwU4f7K3OrvvOsxcs4LHhkS1AsIwW/FCvYX0LeqBFo
-----END PRIVATE KEY-----`

func TestAddress(t *testing.T)  {
	priKey, err := cryptoutil.DecodePriv([]byte(testPriv2))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	priv, err := cryptoutil.UnMarshalPrivateKey(privByte)
	pub := priv.Public().(*ecdsa.PublicKey)

	address, err := cryptoutil.PublicKeyToAddress(pub)
	fmt.Println(address)
}
