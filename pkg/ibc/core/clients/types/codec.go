package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/tendermint/go-amino"
)


var IBCClientCodec *codec.Codec

func init(){
	IBCClientCodec = codec.New()
	RegisterCodec(IBCClientCodec)
	codec.RegisterCrypto(IBCClientCodec)

	IBCClientCodec.Seal()
}
func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterInterface((*exported.ClientState)(nil), &amino.InterfaceOptions{AlwaysDisambiguate: true})
	//cdc.RegisterConcrete((*exported.ClientState)(nil), "ibc/clientstate", nil)
	cdc.RegisterInterface((*exported.ConsensusState)(nil), nil)
	cdc.RegisterConcrete(&MsgCreateClient{}, "ibcclient/MsgCreateClient", nil)

	cdc.RegisterConcrete(&GenesisState{}, "ibcclient/GenesisState", nil)
	cdc.RegisterConcrete(&ConsensusStateWithHeight{}, "ibcclient/ConsensusStateWithHeight", nil)
	cdc.RegisterConcrete(&ClientConsensusStates{}, "ibcclient/ClientConsensusStates", nil)
	cdc.RegisterConcrete(&IdentifiedClientState{}, "ibcclient/IdentifiedClientState", nil)

	cdc.RegisterConcrete(&QueryClientStateResponse{}, "ibcclient/QueryClientStateResponse", nil)
	cdc.RegisterConcrete(&QueryClientStatesResponse{}, "ibcclient/QueryClientStatesResponse", nil)

}

