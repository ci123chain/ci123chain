package keeper

import (
	store2 "github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"
	"strings"
)


func (q Keeper)ClientState(ctx sdk.Context, req abci.RequestQuery) ([]byte, error) {
	var reqClientState types.QueryClientStateRequest
	if err := types.IBCClientCodec.UnmarshalJSON(req.Data, &reqClientState); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := host.ClientIdentifierValidator(reqClientState.ClientId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientState, found := q.GetClientState(ctx, reqClientState.ClientId)
	if !found {
		return nil, status.Error(
			codes.NotFound,
			sdkerrors.Wrap(types.ErrClientNotFound, reqClientState.ClientId).Error(),
		)	}
	proofHeight := types.GetSelfHeight(ctx)
	resp := types.QueryClientStateResponse{
		ClientState: clientState,
		ProofHeight: proofHeight,
	}
	return types.IBCClientCodec.MustMarshalJSON(resp), nil
}


func (q Keeper)ClientStates(ctx sdk.Context, req abci.RequestQuery) ([]byte, error) {
	var reqClientState types.QueryClientStatesRequest
	if err := types.IBCClientCodec.UnmarshalJSON(req.Data, &reqClientState); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientStates := types.IdentifiedClientStates{}
	store := store2.NewPrefixStore(ctx.KVStore(q.storeKey), host.KeyClientStorePrefix)

	pageRes, err := pagination.Paginate(store, reqClientState.Pagination, func(key, value []byte) error {
		keySplit := strings.Split(string(key), "/")
		if keySplit[len(keySplit)-1] != "clientState" {
			return nil
		}

		clientState, err := q.UnmarshalClientState(value)
		if err != nil {
			return err
		}

		clientID := keySplit[1]
		if err := host.ClientIdentifierValidator(clientID); err != nil {
			return err
		}

		identifiedClient := types.NewIdentifiedClientState(clientID, clientState)
		clientStates = append(clientStates, identifiedClient)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Sort(clientStates)

	resp := types.QueryClientStatesResponse{
		ClientStates: clientStates,
		Pagination:   pageRes,
	}
	return types.MustMarshalClientStateResp(types.IBCClientCodec, resp), nil
}


// ConsensusState implements the Query/ConsensusState method
func (q Keeper) ConsensusState(ctx sdk.Context,  req abci.RequestQuery ) ([]byte, error) {
	var reqConsensusState types.QueryConsensusStateRequest
	if err := types.IBCClientCodec.UnmarshalJSON(req.Data, &reqConsensusState); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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

