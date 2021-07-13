package types

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math"
)

type ValSet struct {
	Nonce uint64
	Members []ValSetMember
}

func (v ValSet) FilterEmptyAddress() ([]string, []uint64) {
	var addresses []string
	var powers []uint64
	for _, val := range v.Members {
		if val.EthAddress == nil {
			addresses = append(addresses, "")
		} else {
			addresses = append(addresses, val.EthAddress.String())
		}
		powers = append(powers, val.Power)
	}

	return addresses, powers
}

func (v ValSet) OrderSigs(signedMsg []byte, sigs []Confirm) ([]GravitySignature, error) {
	status, err := v.getSignatureStatus(signedMsg, sigs)
	if err != nil {
		return nil, err
	}

	if gravityPowerToPercent(status.PowerOfGoodSigs) >= float32(66) {
		return status.OrderedSignatures, nil
	} else {
		return nil, errors.New("Not enough power of good sigs")
	}
}

func (v ValSet) getSignatureStatus(signedMsg []byte, sigs []Confirm) (SignatureStatus, error) {
	if len(sigs) == 0 {
		return SignatureStatus{}, errors.New("No signatures")
	}

	var out []GravitySignature
	signatureMap := make(map[common.Address]Confirm)
	for _, sig := range sigs {
		signatureMap[sig.GetEthAddress()] = sig
	}
	var power_of_good_sigs uint64 = 0;
	var power_of_unset_keys uint64 = 0;
	var number_of_unset_key_validators uint64 = 0;
	var power_of_nonvoters uint64 = 0;
	var number_of_nonvoters uint64 = 0;
	var power_of_invalid_signers uint64 = 0;
	var number_of_invalid_signers uint64 = 0;
	number_validators := uint64(len(v.Members))

	for _, member := range v.Members {
		if member.EthAddress != nil {
			ethAddress := *member.EthAddress
			if sig := signatureMap[ethAddress]; sig != nil {
				if sig.GetEthAddress().String() != ethAddress.String() {
					panic("sigs.EthAddress not equal to member.EthAddress")
				}
				recoverKey, err := sig.GetSignature().Ecreover(signedMsg)
				if err != nil {
					panic("recover sig failed")
				}
				if recoverKey.String() == sig.GetEthAddress().String() {
					out = append(out, GravitySignature{
						Power:      member.Power,
						EthAddress: sig.GetEthAddress(),
						V:          sig.GetSignature().V,
						R:          sig.GetSignature().R,
						S:          sig.GetSignature().S,
					})
					power_of_good_sigs += member.Power
				} else {
					out = append(out, GravitySignature{
						Power:      member.Power,
						EthAddress: ethAddress,
						V:          nil,
						R:          nil,
						S:          nil,
					})
					power_of_invalid_signers += member.Power
					number_of_invalid_signers += 1
				}
			} else {
				out = append(out, GravitySignature{
					Power:      member.Power,
					EthAddress: ethAddress,
					V:          nil,
					R:          nil,
					S:          nil,
				})
				power_of_nonvoters += member.Power
				number_of_nonvoters += 1
			}
		} else {
			out = append(out, GravitySignature{
				Power:      member.Power,
				EthAddress: common.Address{},
				V:          nil,
				R:          nil,
				S:          nil,
			})
			power_of_unset_keys += member.Power
			number_of_unset_key_validators += 1
		}
	}

	return SignatureStatus{
		OrderedSignatures:          out,
		PowerOfGoodSigs:            power_of_good_sigs,
		PowerOfUnsetKeys:           power_of_unset_keys,
		NumberOfUnsetKeyValidators: number_of_unset_key_validators,
		PowerOfNonvoters:           power_of_nonvoters,
		NumberOfNonvoters:          number_of_nonvoters,
		PowerOfInvalidSigners:      power_of_invalid_signers,
		NumberOfInvalidSigners:     number_of_invalid_signers,
		NumValidators:              number_validators,
	}, nil
}

type ValSetMember struct {
	Power uint64
	EthAddress *common.Address
}

type ValsetConfirmResponse struct {
	Nonce        uint64
	Orchestrator common.Address
	EthAddress   common.Address
	Signature    EthSignature
}

func (vcr ValsetConfirmResponse) GetEthAddress() common.Address {
	return vcr.Orchestrator
}

func (vcr ValsetConfirmResponse) GetSignature() EthSignature {
	return vcr.Signature
}

type SignatureStatus struct {
	OrderedSignatures []GravitySignature
	PowerOfGoodSigs uint64
	PowerOfUnsetKeys uint64
	NumberOfUnsetKeyValidators uint64
	PowerOfNonvoters uint64
	NumberOfNonvoters uint64
	PowerOfInvalidSigners uint64
	NumberOfInvalidSigners uint64
	NumValidators uint64
}

func gravityPowerToPercent(input uint64) float32 {
	return float32(input) / float32(math.MaxUint32) * float32(100)
}
