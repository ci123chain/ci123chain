package evmtypes

import (
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"
)

// MarshalBigInt marshalls big int into text string for consistent encoding
func MarshalBigInt(i *big.Int) (string, error) {
	bz, err := i.MarshalText()
	if err != nil {
		return "", err
	}
	return string(bz), nil
}

// MustMarshalBigInt marshalls big int into text string for consistent encoding.
// It panics if an error is encountered.
func MustMarshalBigInt(i *big.Int) string {
	str, err := MarshalBigInt(i)
	if err != nil {
		panic(err)
	}
	return str
}

// UnmarshalBigInt unmarshalls string from *big.Int
func UnmarshalBigInt(s string) (*big.Int, error) {
	ret := new(big.Int)
	err := ret.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// MustUnmarshalBigInt unmarshalls string from *big.Int.
// It panics if an error is encountered.
func MustUnmarshalBigInt(s string) *big.Int {
	ret, err := UnmarshalBigInt(s)
	if err != nil {
		panic(err)
	}
	return ret
}

// recoverEthSig recovers a signature according to the Ethereum specification and
// returns the sender or an error.
//
// Ref: Ethereum Yellow Paper (BYZANTIUM VERSION 69351d5) Appendix F
// nolint: gocritic
func recoverEthSig(R, S, Vb *big.Int, sigHash ethcmn.Hash) (ethcmn.Address, error) {
	if Vb.BitLen() > 8 {
		return ethcmn.Address{}, errors.New("invalid signature")
	}

	V := byte(Vb.Uint64() - 27)
	if !ethcrypto.ValidateSignatureValues(V, R, S, true) {
		return ethcmn.Address{}, errors.New("invalid signature")
	}

	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)

	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V

	// recover the public key from the signature
	pub, err := ethcrypto.Ecrecover(sigHash[:], sig)
	if err != nil {
		return ethcmn.Address{}, err
	}

	if len(pub) == 0 || pub[0] != 4 {
		return ethcmn.Address{}, errors.New("invalid public key")
	}

	var addr ethcmn.Address
	copy(addr[:], ethcrypto.Keccak256(pub[1:])[12:])

	return addr, nil
}

func rlpHash(x interface{}) (hash ethcmn.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	_ = rlp.Encode(hasher, x)
	_ = hasher.Sum(hash[:0])

	return hash
}

// ResultData represents the data returned in an sdk.Result
type ResultData struct {
	ContractAddress ethcmn.Address  `json:"contract_address"`
	Bloom           ethtypes.Bloom  `json:"bloom"`
	Logs            []*ethtypes.Log `json:"logs"`
	Ret             []byte          `json:"ret"`
	TxHash          ethcmn.Hash     `json:"tx_hash"`
}

// String implements fmt.Stringer interface.
func (rd ResultData) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ResultData:
	ContractAddress: %s
	Bloom: %s
	Logs: %v
	Ret: %v
	TxHash: %s
`, rd.ContractAddress.String(), rd.Bloom.Big().String(), rd.Logs, rd.Ret, rd.TxHash.String()))
}

// EncodeResultData takes all of the necessary data from the EVM execution
// and returns the data as a byte slice encoded with amino
func EncodeResultData(data ResultData) ([]byte, error) {
	return ModuleCdc.MarshalBinaryLengthPrefixed(data)
}

// DecodeResultData decodes an amino-encoded byte slice into ResultData
func DecodeResultData(in []byte) (ResultData, error) {
	var data ResultData
	err := ModuleCdc.UnmarshalBinaryLengthPrefixed(in, &data)
	if err != nil {
		return ResultData{}, err
	}
	return data, nil
}

// ----------------------------------------------------------------------------
// Auxiliary

// TxDecoder returns an sdk.TxDecoder that can decode both auth.StdTx and
// MsgEthereumTx transactions.
func TxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var tx sdk.Tx

		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrParams, "empty tx")
		}

		// sdk.Tx is an interface. The concrete message types
		// are registered by MakeTxCodec
		// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal failed: %v", err.Error()))
		}

		return tx, nil
	}
}