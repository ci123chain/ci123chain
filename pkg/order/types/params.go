package types

import (
	"time"
)

const SleepTime = 1 * time.Second
const StateDone = "Done"
const StateCommitting = "Committing"
const StateInit = "Init"
const OrderBookKey = "OrderBook"
const OpADD = "ADD"
const NoOrderBookErr = "No OrderBook"

var (
	KeyOrderBook	= []byte("OrderBook")
	StoreKey		= "order"
)

type OrderBook struct {
	Lists 	[]Lists 	`json:"lists"`

	Current	Current 	`json:"current"`

	Actions	[]Actions 	`json:"actions"`
}

type Lists struct {
	Name 	string 	`json:"name"`
	Height	int64	`json:"height"`
}

type Current struct {
	Index	int		`json:"index"`
	State	string	`json:"state"`
}

type Actions struct {
	Type	string	`json:"types"`
	Height	int64	`json:"height"`
	Name	string	`json:"name"`
}

type Params struct {
	OrderBook	OrderBook	`json:"orderBook"`
}

func DefaultParams() Params {
	var lists []Lists

	p1 := &Lists{
		Name:   "",
		Height: 0,
	}
	//p2 := &keeper.Lists{
	//	Name:   "Shard2",
	//	Height: 0,
	//}
	//lists = append(lists, *p1, *p2)
	lists = append(lists, *p1)

	current := Current{
		Index: 0,
		State: StateInit,
	}

	var actions []Actions
	orderbook := OrderBook{
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

