package cryptosuit

type SignIdentity interface {
	// 签名
	Sign(msg []byte, identity []byte) ([]byte, error)

	Verifier(msg []byte, signature []byte, pub []byte, address []byte) (bool, error)

	GetPubKey([]byte) ([]byte, error)
}

const (
	ETHSignType uint8 = 0 + iota
	FabSignType
)

var signMap = map[uint8]SignIdentity{
	ETHSignType: ethsignimp{},
	FabSignType: fabsignimp{},
}

func NewFabSignIdentity() fabsignimp {
	return fabsignimp{}
}

func NewETHSignIdentity() ethsignimp {
	return ethsignimp{}
}

//
//func GetSignIdentity(t uint8) SignIdentity {
//	if sid, ok := signMap[t]; ok {
//		return sid
//	}
//	panic("unsupported cryptosuite")
//}
