package types

import (
	"bytes"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	ETHEREUM_SALT = "\x19Ethereum Signed Message:\n32"
	B32 = "bytes32"
	U256 = "uint256"
	U256Ary = "uint256[]"
	Addr = "address"
	AddrAry = "address[]"
)

type Confirm interface {
	GetEthAddress() common.Address
	GetSignature() EthSignature
}

func EncodeValsetConfirm(gravityId string, valset ValSet) []byte {
	ethAddresses, powers := valset.FilterEmptyAddress()

	return utils.RawEncode([]string{B32, B32, U256, AddrAry, U256Ary},
	[]interface{}{[]byte(gravityId), []byte("checkpoint"), valset.Nonce, ethAddresses, powers})
}

func BytesCombine(pBytes ...[]byte) []byte {
	var buffer bytes.Buffer
	for index := 0; index < len(pBytes); index++ {
		buffer.Write(pBytes[index])
	}
	return buffer.Bytes()
}

func GetEthereumMsgHash(msg []byte) []byte {
	msg = crypto.Keccak256(msg)
	saltMsg := append([]uint8(ETHEREUM_SALT), msg...)
	msgHash := crypto.Keccak256(saltMsg)
	return msgHash
}

func EncodeTxBatchConfirmHashed(gravityId string, batch TransactionBatch) []byte {
	msg := EncodeTxBatchConfirm(gravityId, batch)
	x := GetEthereumMsgHash(msg)
	return x[:]
}

func EncodeValsetConfirmHashed(gravityId string, valset ValSet) []byte {
	msg := EncodeValsetConfirm(gravityId, valset)
	x := GetEthereumMsgHash(msg)
	return x[:]
}

func EncodeTxBatchConfirm(gravityId string, batch TransactionBatch) []byte {
	amounts, fees, destinations := batch.GetCheckPointValues()
	return utils.RawEncode([]string{B32, B32, U256Ary, AddrAry, U256Ary, U256, Addr, U256},
	[]interface{}{[]byte(gravityId), []byte("transactionBatch"), amounts, destinations, fees, batch.Nonce, batch.TokenContract.String(), batch.BatchTimeout})
}