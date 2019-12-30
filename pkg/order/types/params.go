package types

import (
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	"github.com/tanhuiya/ci123chain/pkg/params/subspace"
)

var (
	KeyOrderBook	= []byte("OrderBook")
)

type Params struct {
	OrderBook	keeper.OrderBook	`json:"orderBook"`
}

func DefaultParams() Params {
	var turns []keeper.Turns
	turn1 := &keeper.Turns{
		Name:   "Shard1",
		Height: 0,
	}
	turn2 := &keeper.Turns{
		Name:   "Shard2",
		Height: 0,
	}
	turns = append(turns, *turn1, *turn2)
	//turns = append(turns, *turn1)

	current := keeper.Current{
		Index: 0,
		State: keeper.StateInit,
	}

	var event []keeper.Events
	orderbook := keeper.OrderBook{
		Turns: turns,
		Current: current,
		Events: event,
	}
	return Params{
		OrderBook:  orderbook,
	}
}

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{Params: params}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams())
}

func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyOrderBook, &p.OrderBook},
	}
}