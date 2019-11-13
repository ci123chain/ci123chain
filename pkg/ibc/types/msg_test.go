package types

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"testing"
)

var testCert = `-----BEGIN CERTIFICATE-----
MIICGTCCAcCgAwIBAgIRALR/1GXtEud5GQL2CZykkOkwCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
Lm9yZzEuZXhhbXBsZS5jb20wHhcNMTcwNzI4MTQyNzIwWhcNMjcwNzI2MTQyNzIw
WjBbMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
U2FuIEZyYW5jaXNjbzEfMB0GA1UEAwwWVXNlcjFAb3JnMS5leGFtcGxlLmNvbTBZ
MBMGByqGSM49AgEGCCqGSM49AwEHA0IABPIVPS+hdftwDg8+02y1aV5pOnCO9tIn
f60wZMbrt/5N0J8PFZgylBjEuUTxWRsTMpYPAJi8NlEwoJB+/YSs29ujTTBLMA4G
A1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAIIeR0TY+iVFf
mvoEKwaToscEu43ZXSj5fTVJornjxDUtMAoGCCqGSM49BAMCA0cAMEQCID+dZ7H5
AiaiI2BjxnL3/TetJ8iFJYZyWvK//an13WV/AiARBJd/pI5A7KZgQxJhXmmR8bie
XdsmTcdRvJ3TS/6HCA==
-----END CERTIFICATE-----`

func TestUnmarshalSignedTx(t *testing.T)  {
	var obj SignedIBCMsg
	data := `{"signature":"MEUCIQDTUpUwbhX1+sbLLodc5dPaDkoMGbSMBiwHC6x4N08mzAIgDldd0hrEc/vaKxq11gN8ZQnowNp1B9JAsNpjKmQVEhY=","ibc_msg_bytes":"eyJ1bmlxdWVfaWQiOiIwRDVBNkU0NjI2MEJGNTJENzMyNkYxNjZBODZDMDhCQSIsIm9ic2VydmVyX2lkIjoiMTIzNDU2NzgxMjM0NTY3OCIsImJhbmtfYWRkcmVzcyI6IjB4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMCIsImFwcGx5X3RpbWUiOiIyMDE5LTExLTEyVDIwOjI3OjI5LjQ5MDA1MiswODowMCIsInN0YXRlIjoicmVhZHkiLCJmcm9tX2FkZHJlc3MiOiIweDIwNGJDQzQyNTU5RmFmNkRGRTE0ODUyMDhGNzk1MWFhRDgwMEIzMTMiLCJ0b19hZGRyZXNzIjoiMHhEMWExNDk2MjYyN2ZBYzc2OEZlODg1RWViOUZGMDcyNzA2QjU0YzE5IiwiYW1vdW50IjoxMH0="}`
	b := []byte(data)
	err := json.Unmarshal(b, &obj)
	assert.NoError(t, err)
	jsonStr := string(obj.IBCMsgBytes)
	fmt.Println(jsonStr)

	var ibcMsg IBCMsg
	err = json.Unmarshal(obj.IBCMsgBytes, &ibcMsg)
	assert.NoError(t, err)
	fmt.Println(ibcMsg.UniqueID)

	sid := cryptosuit.NewFabSignIdentity()
	pubkey, _ := cryptoutil.DecodePub([]byte(testCert))
	pubketBz := cryptoutil.MarshalPubkey(pubkey)
	valid, err := sid.Verifier(obj.GetSignBytes(), obj.Signature, pubketBz, nil)
	assert.NoError(t, err)
	assert.True(t, valid)
}
