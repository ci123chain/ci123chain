package keeper

import (
	"fmt"
	types2 "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	errors2 "github.com/ci123chain/ci123chain/pkg/ibc/core/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/pkg/errors"
)

func (q Keeper)ClientState(ctx types2.Context, req *types.QueryClientStateRequest) (*types.QueryClientStateResponse, error) {
	if req == nil {
		return nil, errors.New("empty request")
	}

	if err := host.ClientIdentifierValidator(req.ClientId); err != nil {
		return nil, errors.New(fmt.Sprintf("invalid clientID: %s", req.ClientId))
	}

	clientState, found := q.GetClientState(ctx, req.ClientId)
	if !found {
		return nil, errors2.ErrorClientNotFound(errors2.DefaultCodespace, errors.New("clientid: " + req.ClientId))
	}
	proofHeight := types.GetSelfHeight(ctx)
	return &types.QueryClientStateResponse{
		ClientState: clientState,
		ProofHeight: proofHeight,
	}, nil
}