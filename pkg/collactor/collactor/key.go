package collactor

// KeyOutput contains mnemonic and address of key
type KeyOutput struct {
	PrivateKey  string `json:"privatekey" yaml:"privatekey"`
	Address  string `json:"address" yaml:"address"`
}

