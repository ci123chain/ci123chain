package cryptosuit

type ethsignimp struct {

}

// 签名
func (eth ethsignimp)Sign(msg []byte, identity []byte) ([]byte, error) {
	return []byte{}, nil
}

func (eth ethsignimp)Verifier(msg []byte, signature []byte, pub []byte, address []byte) (bool, error) {
	return true, nil
}


func (fab ethsignimp) GetPubKey(privKey []byte) ([]byte, error) {
	return nil, nil
}