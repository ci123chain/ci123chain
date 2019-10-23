package cryptosuit

type SignIdentity interface {
	Sign(msg []byte, priv []byte)
}