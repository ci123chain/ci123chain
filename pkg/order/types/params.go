package types

import (
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	"github.com/tanhuiya/ci123chain/pkg/params/subspace"
)

var (
	KeyOrderBook	= []byte("OrderBook")
	StoreKey		= "order"
)

type Params struct {
	OrderBook	keeper.OrderBook	`json:"orderBook"`
}

func DefaultParams() Params {
	var lists []keeper.Lists

	p1 := &keeper.Lists{
		Name:   "",
		Height: 0,
	}
	//p2 := &keeper.Lists{
	//	Name:   "Shard2",
	//	Height: 0,
	//}
	//lists = append(lists, *p1, *p2)
	lists = append(lists, *p1)

	current := keeper.Current{
		Index: 0,
		State: keeper.StateInit,
	}

	var actions []keeper.Actions
	orderbook := keeper.OrderBook{
		Lists: lists,
		Current: current,
		Actions: actions,
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