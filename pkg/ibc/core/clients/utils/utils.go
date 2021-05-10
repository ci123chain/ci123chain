package utils

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	ibc "github.com/ci123chain/ci123chain/pkg/ibc"
	ibcclient "github.com/ci123chain/ci123chain/pkg/ibc/core/client"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	host "github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	coretypes "github.com/ci123chain/ci123chain/pkg/ibc/core/types"
	sdkerrors "github.com/pkg/errors"
)



// QueryClientStateABCI queries the store to get the light client state and a merkle proof.
func QueryClientStateABCI2(
	clientCtx context.Context, clientID string,
) (*types.QueryClientStateResponse, error) {
	path := "/custom/" + ibc.ModuleName + "/" + coretypes.QueryClientState

	req := &clienttypes.QueryClientStateRequest{
		ClientId: clientID,
	}
	key := clientCtx.Cdc.MustMarshalJSON(req)
	value, _, err := ibcclient.QueryABCI(clientCtx, path, key, false)
	if err != nil {
		return nil, err
	}

	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(types.ErrClientNotFound, "clients not found")
	}

	var resp types.QueryClientStateResponse
	err = clientCtx.Cdc.UnmarshalJSON(value, &resp)
	//clientStatesResp, err := types.UnmarshalClientStateResp(types.IBCClientCodec, value)
	return &resp, err
}



// QueryClientStateABCI queries the store to get the light client state and a merkle proof.
func QueryClientStateABCI(
	clientCtx context.Context, clientID string,
) (*types.QueryClientStateResponse, error) {
	key := host.FullClientStateKey(clientID)

	value, proofBz, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(types.ErrClientNotFound, clientID)
	}


	clientState, err := types.UnmarshalClientState(clientCtx.Cdc, value)
	if err != nil {
		return nil, err
	}

	anyClientState, err := types.PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	clientStateRes := types.NewQueryClientStateResponse(anyClientState, proofBz, proofHeight)
	return clientStateRes, nil
}

// QueryConsensusStateABCI queries the store to get the consensus state of a light client and a
// merkle proof of its existence or non-existence.
func QueryConsensusStateABCI(
	clientCtx context.Context, clientID string, height exported.Height,
) (*types.QueryConsensusStateResponse, error) {
	key := host.FullConsensusStateKey(clientID, height)

	value, proofBz, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	// check if consensus state exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(types.ErrConsensusStateNotFound, clientID)
	}

	cs, err := types.UnmarshalConsensusState(clientCtx.Cdc, value)
	if err != nil {
		return nil, err
	}

	anyConsensusState, err := types.PackConsensusState(cs)
	if err != nil {
		return nil, err
	}

	return types.NewQueryConsensusStateResponse(anyConsensusState, proofBz, proofHeight), nil
}


// QueryClientStateABCI queries the store to get the light client state and a merkle proof.
func QueryClientStatesABCI(
	clientCtx context.Context, offset, limit uint64,
) (*types.QueryClientStatesResponse, error) {
	path := "/custom/" + ibc.ModuleName + "/" + coretypes.QueryClientStates

	req := &clienttypes.QueryClientStatesRequest{
		Pagination: &pagination.PageRequest{
			Key:        []byte(""),
			Offset:     offset,
			Limit:      limit,
			CountTotal: true,
		},
	}
	key := clientCtx.Cdc.MustMarshalJSON(req)
	value, _, err := ibcclient.QueryABCI(clientCtx, path, key, false)
	if err != nil {
		return nil, err
	}

	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(types.ErrClientNotFound, "clients not found")
	}


	clientStatesResp, err := types.UnmarshalClientStateResp(clientCtx.Cdc, value)
	return &clientStatesResp, err
}

