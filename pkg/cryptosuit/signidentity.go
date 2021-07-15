package cryptosuit

type SignIdentity interface {
	// 签名
	Sign(msg []byte, identity []byte) ([]byte, error)

	Verifier(msg []byte, signature []byte, pub []byte, address []byte) (bool, error)

	GetPubKey([]byte) ([]byte, error)
}

const (
	ETHSignType uint8 = 0 + iota
)

var signMap = map[uint8]SignIdentity{
	ETHSignType: ethsignimp{},
}


func NewETHSignIdentity() ethsignimp {
	return ethsignimp{}
}

