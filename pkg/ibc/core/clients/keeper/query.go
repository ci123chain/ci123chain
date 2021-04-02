package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	errors2 "github.com/ci123chain/ci123chain/pkg/ibc/core/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (q Keeper)ClientState(ctx sdk.Context, req abci.RequestQuery) ([]byte, error) {
	var reqClientState types.QueryClientStateRequest
	if err := types.IBCClientCodec.UnmarshalJSON(req.Data, &reqClientState); err != nil {
		return nil, err
	}

	if err := host.ClientIdentifierValidator(reqClientState.ClientId); err != nil {
		return nil, errors.New(fmt.Sprintf("invalid clientID: %s", reqClientState.ClientId))
	}

	clientState, found := q.GetClientState(ctx, reqClientState.ClientId)
	if !found {
		return nil, errors2.ErrorClientNotFound(errors2.DefaultCodespace, errors.New("clientid: " + reqClientState.ClientId))
	}
	proofHeight := types.GetSelfHeight(ctx)
	resp := types.QueryClientStateResponse{
		ClientState: clientState,
		ProofHeight: proofHeight,
	}
	return types.IBCClientCodec.MustMarshalJSON(resp), nil
}


// ConsensusState implements the Query/ConsensusState method
func (q Keeper) ConsensusState(ctx sdk.Context,  req abci.RequestQuery ) ([]byte, error) {
	var reqConsensusState types.QueryConsensusStateRequest
	if err := types.IBCClientCodec.UnmarshalJSON(req.Data, &reqConsensusState); err != nil {
		return nil, err
	}

	if err := host.ClientIdentifierValidator(reqConsensusState.ClientId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		consensusState exported.ConsensusState
		found          bool
	)

	height := types.NewHeight(reqConsensusState.RevisionNumber, reqConsensusState.RevisionHeight)
	if reqConsensusState.LatestHeight {
		consensusState, found = q.GetLatestClientConsensusState(ctx, reqConsensusState.ClientId)
	} else {
		if reqConsensusState.RevisionHeight == 0 {
			return nil, status.Error(codes.InvalidArgument, "consensus state height cannot be 0")
		}

		consensusState, found = q.GetClientConsensusState(ctx, reqConsensusState.ClientId, height)
	}

	if !found {
		return nil, errors.Wrapf(types.ErrConsensusStateNotFound, "client-id: %s, height: %s", reqConsensusState.ClientId, height)
	}

	proofHeight := types.GetSelfHeight(ctx)
	resp := types.QueryConsensusStateResponse{
		ConsensusState: consensusState,
		ProofHeight:    proofHeight,
	}
	return types.IBCClientCodec.MustMarshalJSON(resp), nil
}