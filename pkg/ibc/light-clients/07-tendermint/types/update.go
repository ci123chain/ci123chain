package types

import (
	"bytes"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/tendermint/tendermint/light"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
)

func (cs ClientState) CheckHeaderAndUpdateState(ctx sdk.Context, cdc *codec.Codec, clientStore sdk.KVStore,
	header exported.Header) (exported.ClientState, exported.ConsensusState, error) {
	tmHeader, ok := header.(*Header)
	if !ok {
		return nil, nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader, "expected type %T, got %T", &Header{}, header,
		)
	}

	// get consensus state from clientStore
	tmConsState, err := GetConsensusState(clientStore, cdc, tmHeader.TrustedHeight)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(
			err, "could not get consensus state from clientstore at TrustedHeight: %s", tmHeader.TrustedHeight,
		)
	}

	if err := checkValidity(&cs, tmConsState, tmHeader, ctx.BlockHeader().Time); err != nil {
		return nil, nil, err
	}

	newClientState, consensusState := update(ctx, clientStore, &cs, tmHeader)
	return newClientState, consensusState, nil
}


// checkValidity checks if the Tendermint header is valid.
// CONTRACT: consState.Height == header.TrustedHeight
func checkValidity(
	clientState *ClientState, consState *ConsensusState,
	header *Header, currentTimestamp time.Time,
) error {
	if err := checkTrustedHeader(header, consState); err != nil {
		return err
	}

	// UpdateClient only accepts updates with a header at the same revision
	// as the trusted consensus state
	if header.GetHeight().GetRevisionNumber() != header.TrustedHeight.RevisionNumber {
		return sdkerrors.Wrapf(
			ErrInvalidHeaderHeight,
			"header height revision %d does not match trusted header revision %d",
			header.GetHeight().GetRevisionNumber(), header.TrustedHeight.RevisionNumber,
		)
	}

	//tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	//if err != nil {
	//	return sdkerrors.Wrap(err, "trusted validator set in not tendermint validator set type")
	//}
	//
	//tmSignedHeader, err := tmtypes.SignedHeaderFromProto(header.SignedHeader)
	//if err != nil {
	//	return sdkerrors.Wrap(err, "signed header in not tendermint signed header type")
	//}
	//
	//tmValidatorSet, err := tmtypes.ValidatorSetFromProto(header.ValidatorSet)
	//if err != nil {
	//	return sdkerrors.Wrap(err, "validator set in not tendermint validator set type")
	//}

	// assert header height is newer than consensus state
	if header.GetHeight().LTE(header.TrustedHeight) {
		return sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader,
			"header height ≤ consensus state height (%s ≤ %s)", header.GetHeight(), header.TrustedHeight,
		)
	}

	chainID := clientState.GetChainID()
	// If chainID is in revision format, then set revision number of chainID with the revision number
	// of the header we are verifying
	// This is useful if the update is at a previous revision rather than an update to the latest revision
	// of the client.
	// The chainID must be set correctly for the previous revision before attempting verification.
	// Updates for previous revisions are not supported if the chainID is not in revision format.
	if clienttypes.IsRevisionFormat(chainID) {
		chainID, _ = clienttypes.SetRevisionNumber(chainID, header.GetHeight().GetRevisionNumber())
	}

	// Construct a trusted header using the fields in consensus state
	// Only Height, Time, and NextValidatorsHash are necessary for verification
	trustedHeader := tmtypes.Header{
		ChainID:            chainID,
		Height:             int64(header.TrustedHeight.RevisionHeight),
		Time:               consState.Timestamp,
		NextValidatorsHash: consState.NextValidatorsHash,
	}
	signedHeader := tmtypes.SignedHeader{
		Header: &trustedHeader,
	}

	// Verify next header with the passed-in trustedVals
	// - asserts trusting period not passed
	// - assert header timestamp is not past the trusting period
	// - assert header timestamp is past latest stored consensus state timestamp
	// - assert that a TrustLevel proportion of TrustedValidators signed new Commit
	err := light.Verify(
		&signedHeader,
		header.TrustedValidators, header.SignedHeader, header.ValidatorSet,
		clientState.TrustingPeriod, currentTimestamp, clientState.MaxClockDrift, clientState.TrustLevel.ToTendermint(),
	)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to verify header")
	}
	return nil
}

// update the consensus state from a new header and set processed time metadata
func update(ctx sdk.Context, clientStore sdk.KVStore, clientState *ClientState, header *Header) (*ClientState, *ConsensusState) {
	height := header.GetHeight().(clienttypes.Height)
	if height.GT(clientState.LatestHeight) {
		clientState.LatestHeight = height
	}
	consensusState := &ConsensusState{
		Timestamp:          header.GetTime(),
		Root:               commitmenttypes.NewMerkleRoot(header.Header.AppHash),
		NextValidatorsHash: header.Header.NextValidatorsHash,
	}

	// set context time as processed time as this is state internal to tendermint client logic.
	// client state and consensus state will be set by client keeper
	SetProcessedTime(clientStore, header.GetHeight(), uint64(ctx.BlockHeader().Time.UnixNano()))

	return clientState, consensusState
}


// checkTrustedHeader checks that consensus state matches trusted fields of Header
func checkTrustedHeader(header *Header, consState *ConsensusState) error {
	//tmTrustedValidators, err := tmtypes.ValidatorSetFromProto(header.TrustedValidators)
	//if err != nil {
	//	return sdkerrors.Wrap(err, "trusted validator set in not tendermint validator set type")
	//}

	// assert that trustedVals is NextValidators of last trusted header
	// to do this, we check that trustedVals.Hash() == consState.NextValidatorsHash
	tvalHash := header.TrustedValidators.Hash()
	if !bytes.Equal(consState.NextValidatorsHash, tvalHash) {
		return sdkerrors.Wrapf(
			ErrInvalidValidatorSet,
			"trusted validators %s, does not hash to latest trusted validators. Expected: %X, got: %X",
			header.TrustedValidators, consState.NextValidatorsHash, tvalHash,
		)
	}
	return nil
}
