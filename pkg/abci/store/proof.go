package store

import (
	ics23 "github.com/confio/ics23/go"
	tmmerkle "github.com/tendermint/tendermint/proto/tendermint/crypto"
)



const (
	ProofOpIAVLCommitment         = "ics23:iavl"
	ProofOpSimpleMerkleCommitment = "ics23:simple"
)

type CommitmentOp struct {
	Type string
	Spec *ics23.ProofSpec
	Key  []byte
	Proof *ics23.CommitmentProof
}

// ProofOp implements ProofOperator interface and converts a CommitmentOp
// into a merkle.ProofOp format that can later be decoded by CommitmentOpDecoder
// back into a CommitmentOp for proof verification
func (op CommitmentOp) ProofOp() tmmerkle.ProofOp {
	bz, err := op.Proof.Marshal()
	if err != nil {
		panic(err.Error())
	}
	return tmmerkle.ProofOp{
		Type: op.Type,
		Key:  op.Key,
		Data: bz,
	}
}

func NewIavlCommitmentOp(key []byte, proof *ics23.CommitmentProof) CommitmentOp {
	return CommitmentOp{
		Type: ProofOpIAVLCommitment,
		Spec: ics23.IavlSpec,
		Key: key,
		Proof: proof,
	}
}

