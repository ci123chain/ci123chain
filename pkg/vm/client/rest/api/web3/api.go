package web3

// PublicWeb3API is the web3_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicWeb3API struct{}

// New creates an instance of the Web3 API.
func NewAPI() *PublicWeb3API {
	return &PublicWeb3API{}
}

// ClientVersion returns the clients version in the Web3 user agent format.
func (PublicWeb3API) ClientVersion() string {
	return ClientVersion()
}
