package types

import (
	"bytes"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	ics23 "github.com/confio/ics23/go"
	"github.com/pkg/errors"
	"net/url"
)

var _ exported.Root = (*MerkleRoot)(nil)
// MerkleRoot defines a merkle root hash.
// In the Cosmos SDK, the AppHash of a block header becomes the root.
type MerkleRoot struct {
	Hash []byte `json:"hash,omitempty"`
}

// NewMerkleRoot constructs a new MerkleRoot
func NewMerkleRoot(hash []byte) MerkleRoot {
	return MerkleRoot{
		Hash: hash,
	}
}

// GetHash implements RootI interface
func (mr MerkleRoot) GetHash() []byte {
	return mr.Hash
}

// Empty returns true if the root is empty
func (mr MerkleRoot) Empty() bool {
	return len(mr.GetHash()) == 0
}


// MerklePath is the path used to verify commitment proofs, which can bean
// arbitrary structured object (defined by a commitment type).
// MerklePath is represented from root-to-leaf
type MerklePath struct {
	KeyPath []string `json:"key_path,omitempty" yaml:"key_path"`
}


var _ exported.Path = (*MerklePath)(nil)

// NewMerklePath creates a new MerklePath instance
// The keys must be passed in from root-to-leaf order
func NewMerklePath(keyPath ...string) MerklePath {
	return MerklePath{
		KeyPath: keyPath,
	}
}

// String implements fmt.Stringer.
// This represents the path in the same way the tendermint KeyPath will
// represent a key path. The backslashes partition the key path into
// the respective stores they belong to.
func (mp MerklePath) String() string {
	pathStr := ""
	for _, k := range mp.KeyPath {
		pathStr += "/" + url.PathEscape(k)
	}
	return pathStr
}

// Pretty returns the unescaped path of the URL string.
// This function will unescape any backslash within a particular store key.
// This makes the keypath more human-readable while removing information
// about the exact partitions in the key path.
func (mp MerklePath) Pretty() string {
	path, err := url.PathUnescape(mp.String())
	if err != nil {
		panic(err)
	}
	return path
}

// GetKey will return a byte representation of the key
// after URL escaping the key element
func (mp MerklePath) GetKey(i uint64) ([]byte, error) {
	if i >= uint64(len(mp.KeyPath)) {
		return nil, fmt.Errorf("index out of range. %d (index) >= %d (len)", i, len(mp.KeyPath))
	}
	key, err := url.PathUnescape(mp.KeyPath[i])
	if err != nil {
		return nil, err
	}
	return []byte(key), nil
}

// Empty returns true if the path is empty
func (mp MerklePath) Empty() bool {
	return len(mp.KeyPath) == 0
}

func ApplyPrefix(prefix exported.Prefix, path MerklePath) (MerklePath, error) {
	if prefix == nil || prefix.Empty() {
		return MerklePath{}, errors.New("prefix can't be empty")
	}
	return NewMerklePath(append([]string{string(prefix.Bytes())}, path.KeyPath...)...), nil
}


// merkle Proof implements
type MerkleProof struct {
	Proofs []*ics23.CommitmentProof `json:"proofs,omitempty"`
}
// Empty returns true if the root is empty
func (proof *MerkleProof) Empty() bool {
	return proof == nil || len(proof.Proofs) == 0
}

// ValidateBasic checks if the proof is empty.
func (proof MerkleProof) ValidateBasic() error {
	if proof.Empty() {
		return errors.New("Invalid proof")
	}
	return nil
}

func (proof MerkleProof) VerifyMembership(specs []*ics23.ProofSpec,
	root exported.Root, path exported.Path, value []byte) error {
	if err := proof.validateVerificationArgs(specs, root); err != nil {
		return err
	}
	mpath, ok := path.(MerklePath)
	if !ok {
		return errors.Errorf("path %v is not of type MerklePath", path)
	}
	if len(mpath.KeyPath) != len(specs) {
		return errors.Errorf("path length %d not same as proof %d", len(mpath.KeyPath), len(specs))
	}
	if len(value) == 0 {
		return errors.Errorf("empty value in membership proof")
	}
	// Since every proof in chain is a membership proof we can use verifyChainedMembershipProof from index 0
	// to validate entire proof
	if err := verifyChainedMembershipProof(root.GetHash(), specs, proof.Proofs, mpath, value, 0); err != nil {
		return err
	}
	return nil
}

func (proof MerkleProof) validateVerificationArgs(specs []*ics23.ProofSpec,
	root exported.Root) error {
	if proof.Empty() {
		return errors.New("proof cannot be empty")
	}
	if root == nil || root.Empty() {
		return errors.New("root cannot be empty")
	}
	if len(specs) != len(proof.Proofs) {
		return errors.Errorf("length of specs: %d not equal to length of proof: %d",
			len(specs), len(proof.Proofs))
	}

	for i, spec := range specs {
		if spec == nil {
			return errors.Errorf("spec at position %d is nil", i)
		}
	}
	return nil
}

// MerklePrefix is merkle path prefixed to the key.
// The constructed key from the Path and the key will be append(Path.KeyPath,
// append(Path.KeyPrefix, key...))
type MerklePrefix struct {
	KeyPrefix []byte `json:"key_prefix,omitempty" yaml:"key_prefix"`
}

var _ exported.Prefix = (*MerklePrefix)(nil)

// NewMerklePrefix constructs new MerklePrefix instance
func NewMerklePrefix(keyPrefix []byte) MerklePrefix {
	return MerklePrefix{
		KeyPrefix: keyPrefix,
	}
}

// Bytes returns the key prefix bytes
func (mp MerklePrefix) Bytes() []byte {
	return mp.KeyPrefix
}

// Empty returns true if the prefix is empty
func (mp MerklePrefix) Empty() bool {
	return len(mp.Bytes()) == 0
}

func verifyChainedMembershipProof(root []byte,
	specs []*ics23.ProofSpec,
	proofs []*ics23.CommitmentProof,
	keys MerklePath, value []byte, index int) error {
	var (
		subroot []byte
		err error
	)
	subroot = value
	for i := index; i < len(proofs); i++ {
		switch proofs[i].Proof.(type) {
			case *ics23.CommitmentProof_Exist:
				subroot, err = proofs[i].Calculate()
				if err != nil {
					return errors.Errorf("could not calculate proof root at index %d, merkle tree may be empty. %v", i, err)
				}

				// Since keys are passed in from highest to lowest, we must grab their indices in reverse order
				// from the proofs and specs which are lowest to highest
				key, err := keys.GetKey(uint64(len(keys.KeyPath) - 1 - i))
				if err != nil {
					return errors.Errorf("could not retrieve key bytes for key %s: %v", keys.KeyPath[len(keys.KeyPath)-1-i], err)
				}
				// verify membership of the proof at this index with appropriate key and value
				if ok := ics23.VerifyMembership(specs[i], subroot, proofs[i], key, value); !ok {
					return errors.Errorf("chained membership proof failed to verify membership of value: %X in subroot %X at index %d. Please ensure the path and value are both correct.",
						value, subroot, i)
				}
				// Set value to subroot so that we verify next proof in chain commits to this subroot
				value = subroot

			case *ics23.CommitmentProof_Nonexist:
				return errors.Errorf("chained membership proof contains nonexistence proof at index %d. If this is unexpected, please ensure that proof was queried from the height that contained the value in store and was queried with the correct key.",
					i)
			default:
				return errors.Errorf("expected proof type: %T, got: %T", &ics23.CommitmentProof_Exist{}, proofs[i].Proof)
		}
	}
	// Check that chained proof root equals passed-in root
	if !bytes.Equal(root, subroot) {
		return errors.Errorf("proof did not commit to expected root: %X, got: %X. Please ensure proof was submitted with correct proofHeight and to the correct chain.",
			root, subroot)
	}
	return nil
}
















