package codec

import (
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// amino codec to marshal/unmarshal
type Codec = amino.Codec

func New() *Codec {
	cdc := amino.NewCodec()
	return cdc
}

// Register the go-crypto to the codec
func RegisterCrypto(cdc *Codec) {
	cryptoAmino.RegisterAmino(cdc)
}

//__________________________________________________________________

// generic sealed codec to be used throughout sdk
var Cdc *Codec

func init() {
	cdc := New()
	RegisterCrypto(cdc)
	Cdc = cdc.Seal()
}
