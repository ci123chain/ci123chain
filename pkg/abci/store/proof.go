package store

import ics23 "github.com/confio/ics23/go"


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

func NewIavlCommitmentOp(key []byte, proof *ics23.CommitmentProof) CommitmentOp {
	return CommitmentOp{
		Type: ProofOpIAVLCommitment,
		Spec: ics23.IavlSpec,
		Key: key,
		Proof: proof,
	}
}