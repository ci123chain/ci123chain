package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/pkg/errors"
)

var KeyProcessedTime = []byte("/processedTime")



// ProcessedTimeKey returns the key under which the processed time will be stored in the client store.
func ProcessedTimeKey(height exported.Height) []byte {
	return append(host.ConsensusStateKey(height), KeyProcessedTime...)
}

func SetProcessedTime(clientStore sdk.KVStore, height exported.Height, timeNs uint64)  {
	key := ProcessedTimeKey(height)
	val := sdk.Uint64ToBigEndian(timeNs)
	clientStore.Set(key, val)
}

func GetConsensusState(store sdk.KVStore, cdc *codec.Codec, height exported.Height) (*ConsensusState, error) {
	bz := store.Get(host.ConsensusStateKey(height))
	if bz == nil {
		return nil, errors.Errorf(
			"consensus state does not exist for height %s", height,
		)
	}
	consensusStateI, err := clienttypes.UnmarshalConsensusState(cdc, bz)
	if err != nil {
		return nil, errors.Errorf("unmarshal error: %v", err)
	}
	consensusState, ok := consensusStateI.(*ConsensusState)

	if !ok {
		return nil, errors.Errorf(
			"invalid consensus type %T, expected %T", consensusState, &ConsensusState{},
		)
	}

	return consensusState, nil
}