package types

import (
	"bytes"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	tmtypes "github.com/tendermint/tendermint/types"
	types2 "github.com/tendermint/tendermint/proto/tendermint/types"
	"time"
)
var _ exported.Header = &Header{}

type Header struct {
	*types2.SignedHeader `protobuf:"bytes,1,opt,name=signed_header,json=signedHeader,proto3,embedded=signed_header" json:"signed_header,omitempty" yaml:"signed_header"`
	ValidatorSet         *types2.ValidatorSet `protobuf:"bytes,2,opt,name=validator_set,json=validatorSet,proto3" json:"validator_set,omitempty" yaml:"validator_set"`
	TrustedHeight        clienttypes.Height         `protobuf:"bytes,3,opt,name=trusted_height,json=trustedHeight,proto3" json:"trusted_height" yaml:"trusted_height"`
	TrustedValidators    *types2.ValidatorSet `protobuf:"bytes,4,opt,name=trusted_validators,json=trustedValidators,proto3" json:"trusted_validators,omitempty" yaml:"trusted_validators"`
}

// ConsensusState returns the updated consensus state associated with the header
func (h Header) ConsensusState() *ConsensusState {
	return &ConsensusState{
		Timestamp:          h.GetTime(),
		Root:               commitmenttypes.NewMerkleRoot(h.Header.GetAppHash()),
		NextValidatorsHash: h.Header.NextValidatorsHash,
	}
}


func (h Header) ClientType() string {
	return exported.Tendermint
}

func (h Header) GetHeight() exported.Height {
	revision := clienttypes.ParseChainID(h.Header.ChainID)
	return clienttypes.NewHeight(revision, uint64(h.Header.Height))
}

func (h Header) ValidateBasic() error {
	if h.SignedHeader == nil {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "tendermint signed header cannot be nil")
	}
	if h.Header == nil {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "tendermint header cannot be nil")
	}
	tmSignedHeader, err := tmtypes.SignedHeaderFromProto(h.SignedHeader)
	if err != nil {
		return sdkerrors.Wrap(err, "header is not a tendermint header")
	}
	if err := tmSignedHeader.ValidateBasic(h.Header.GetChainID()); err != nil {
		return sdkerrors.Wrap(err, "header failed basic validation")
	}

	// TrustedHeight is less than Header for updates
	// and less than or equal to Header for misbehaviour
	if h.TrustedHeight.GT(h.GetHeight()) {
		return sdkerrors.Wrapf(ErrInvalidHeaderHeight, "TrustedHeight %d must be less than or equal to header height %d",
			h.TrustedHeight, h.GetHeight())
	}

	if h.ValidatorSet == nil {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "validator set is nil")
	}
	tmValset, err := tmtypes.ValidatorSetFromProto(h.ValidatorSet)
	if err != nil {
		return sdkerrors.Wrap(err, "validator set is not tendermint validator set")
	}
	if !bytes.Equal(h.Header.ValidatorsHash, tmValset.Hash()) {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "validator set does not match hash")
	}
	return nil
}
// GetTime returns the current block timestamp. It returns a zero time if
// the tendermint header is nil.
// NOTE: the header.Header is checked to be non nil in ValidateBasic.
func (h Header) GetTime() time.Time {
	return h.Header.Time
}