package ethereum_gravity

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
	"github.com/umbracle/go-web3/contract"
	"github.com/umbracle/go-web3/jsonrpc"
	"math/big"
	"testing"
)

func TestAbi(t *testing.T) {
	typ := abi.MustNewType("tuple(address a, uint256 b)")

	type Obj struct {
		A web3.Address
		B *big.Int
	}
	obj := &Obj{
		A: web3.Address{0x1},
		B: big.NewInt(1),
	}

	// Encode
	encoded, err := typ.Encode(obj)
	if err != nil {
		panic(err)
	}

	fmt.Println(encoded)

	// Decode output into a map
	res, err := typ.Decode(encoded)
	if err != nil {
		panic(err)
	}

	// Decode into a struct
	var obj2 Obj
	if err := typ.DecodeStruct(encoded, &obj2); err != nil {
		panic(err)
	}

	fmt.Println(res)
	fmt.Println(obj)
}

var (
	url   = "https://mainnet.infura.io"
	zeroX = web3.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
)

func TestERC20Decimals(t *testing.T) {
	c := &jsonrpc.Client{}
	erc20 := NewERC20(zeroX, c)

	decimals, err := erc20.Decimals()
	assert.NoError(t, err)
	if decimals != 18 {
		t.Fatal("bad")
	}
}

// ERC20 is a solidity contract
type ERC20 struct {
	c *contract.Contract
}

// NewERC20 creates a new instance of the contract at a specific address
func NewERC20(addr web3.Address, provider *jsonrpc.Client) *ERC20 {
	var abiERC20Str = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`
	abiERC20, _ := abi.NewABI(abiERC20Str)
	return &ERC20{c: contract.NewContract(addr, abiERC20, provider)}
}

// Decimals calls the decimals method in the solidity contract
func (a *ERC20) Decimals(block ...web3.BlockNumber) (val0 uint8, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("decimals", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	val0, ok = out["0"].(uint8)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}