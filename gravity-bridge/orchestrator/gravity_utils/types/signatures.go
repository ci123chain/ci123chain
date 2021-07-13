package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type GravitySignature struct {
	Power uint64
	EthAddress common.Address
	V *big.Int
	R *big.Int
	S *big.Int
}

type GravitySignatureArrays  struct {
	Addresses []common.Address
	Powers []uint64
	V []uint8
	R [][]byte
	S [][]byte
}

func ToArrays(signatures []GravitySignature) GravitySignatureArrays {
	var addresses []common.Address
	var powers []uint64
	var v []uint8
	var r [][]byte
	var s [][]byte

	for _, val := range signatures {
		addresses = append(addresses, val.EthAddress)
		powers = append(powers, val.Power)
		if val.V != nil && val.R != nil && val.S != nil{
			v = append(v, val.V.Bytes()[0])
			r = append(r, val.R.Bytes())
			s = append(s, val.S.Bytes())
		} else {
			v = append(v, 0)
			r = append(r, []byte{0})
			s = append(s, []byte{0})
		}
	}

	return GravitySignatureArrays{
		Addresses: addresses,
		Powers:    powers,
		V:         v,
		R:         r,
		S:         s,
	}
}
