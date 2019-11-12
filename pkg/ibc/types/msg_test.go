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
	data := `{"signature":"MEUCIQCkLtPFeI0SCrW8a+kUd5pcKE83Uj45mjW0mjIAItupUgIgRe9eOoZ9aw03GFAB9I3b+N9Ss+BmbViDG/Jan3BF3mE=","ibc_msg_bytes":"eyJ1bmlxdWVfaWQiOiI2OThEODcxQzE1OUYyNUI1MUM0MUFDMDlDNTU1MkZFMyIsIm9ic2VydmVyX2lkIjoiMTIzNDU2NzgxMjM0NTY3OCIsImJhbmtfYWRkcmVzcyI6IjB4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMCIsIkFwcGx5VGltZSI6IjIwMTktMTEtMTJUMTk6NDM6MzkuODE1NTE3KzA4OjAwIiwicmF3IjoiZXlKMGVYQmxJam9pWTJreE1qTmphR0ZwYmk5SlFrTlVjbUZ1YzJabGNpSXNJblpoYkhWbElqcDdJa052YlcxdmJsUjRJanA3SWtOdlpHVWlPakFzSWtaeWIyMGlPaUl3ZURJd05HSkRRelF5TlRVNVJtRm1Oa1JHUlRFME9EVXlNRGhHTnprMU1XRmhSRGd3TUVJek1UTWlMQ0pPYjI1alpTSTZJakUxTnpNMU5USTRPVGczTVRrM09EUXdNREFpTENKSFlYTWlPaUl4SWl3aVVIVmlTMlY1SWpvaVFsQkpWbEJUSzJoa1puUjNSR2M0S3pBeWVURmhWalZ3VDI1RFR6bDBTVzVtTmpCM1drMWljblF2TlU0d1NqaFFSbHBuZVd4Q2FrVjFWVlI0VjFKelZFMXdXVkJCU21rNFRteEZkMjlLUWlzdldWTnpNamx6UFNJc0lsTnBaMjVoZEhWeVpTSTZJazFGVlVOSlVVTXhRM2xMYVhScE5IcExibXBPTlZOT09XeDRlbkU1ZVdFMGVsQkRLekE1VEZsRWRXOVFVME00WVhCQlNXZGlOVzlNYlZselRVcHVjM0V6ZFZJME5FdGxMMVZ6WjBWaVJHaFVTUzlOUjJWbGEyNVdVazFpV0dWalBTSjlMQ0owYjE5aFpHUnlaWE56SWpvaU1IaEVNV0V4TkRrMk1qWXlOMlpCWXpjMk9FWmxPRGcxUldWaU9VWkdNRGN5TnpBMlFqVTBZekU1SWl3aWRXNXBjWFZsWDJsa0lqcHVkV3hzTENKamIybHVJam9pTVRBaWZYMD0iLCJzdGF0ZSI6InJlYWR5In0="}`
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
