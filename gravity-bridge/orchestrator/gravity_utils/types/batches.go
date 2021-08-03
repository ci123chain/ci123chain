package types

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type TransactionBatch struct {
	Nonce uint64
	BatchTimeout uint64
	Transactions []BatchTransaction
	TotalFee Erc20Token
	TokenContract common.Address
}

func (tb TransactionBatch) GetCheckPointValues() ([]*big.Int, []*big.Int, []string) {
	var amounts, fees []*big.Int
	var destinations []string
	for _, transaction := range tb.Transactions {
		amounts = append(amounts, transaction.Erc20Token.Amount.BigInt())
		fees = append(fees, transaction.Erc20Fee.Amount.BigInt())
		destinations = append(destinations, transaction.Destination.String())
	}

	return amounts, fees, destinations
}

type BatchTransaction struct {
	Id uint64
	Sender common.Address
	Destination common.Address
	Erc20Token Erc20Token
	Erc20Fee Erc20Token
}

type BatchConfirmResponse struct {
	Nonce uint64
	Orchestrator common.Address
	TokenContract common.Address
	EthereumSigner common.Address
	EthSignature EthSignature
}

func (bcr BatchConfirmResponse) GetEthAddress() common.Address {
	return bcr.Orchestrator
}

func (bcr BatchConfirmResponse) GetSignature() EthSignature {
	return bcr.EthSignature
}

type EthSignature struct {
	V *big.Int
	R *big.Int
	S *big.Int
}

func (es EthSignature) Ecreover(hash []byte) (common.Address, error) {
	if es.V.BitLen() > 8 {
		return common.Address{}, errors.New("invalid signature")
	}

	V := byte(es.V.Uint64() - 27)
	if !ethcrypto.ValidateSignatureValues(V, es.R, es.S, true) {
		return common.Address{}, errors.New("invalid signature")
	}

	// encode the signature in uncompressed format
	r, s := es.R.Bytes(), es.S.Bytes()
	sig := make([]byte, 65)

	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V

	// recover the public key from the signature
	pub, err := ethcrypto.Ecrecover(hash[:], sig)
	if err != nil {
		return common.Address{}, err
	}

	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}

	var addr common.Address
	copy(addr[:], ethcrypto.Keccak256(pub[1:])[12:])

	return addr, nil
}

func FromBytesToEthSignature(bz []byte) (EthSignature, error) {
	if len(bz) != 65 {
		return EthSignature{}, errors.New("InvalidSignatureLength")
	}

	r := new(big.Int).SetBytes(bz[0:32])
	s := new(big.Int).SetBytes(bz[32:64])
	v := new(big.Int).SetBytes(bz[64:])

	return EthSignature{
		V: v.Add(v, big.NewInt(27)),
		R: r,
		S: s,
	}, nil
}

