package keeper

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"github.com/tanhuiya/ci123chain/pkg/params/subspace"
	"time"
)

var ModuleCdc *codec.Codec
const SleepTime = 1 * time.Second
const StateProcessing = "Processing"
const StateDone = "Done"
const StateInit = "Init"
const OrderBookKey = "OrderBook"
type OrderKeeper struct {
	cdb 		*couchdb.GoCouchDB
	StoreKey	sdk.StoreKey
	paramSubspace subspace.Subspace
}

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
	Type	string	`json:"type"`
	Height	int64	`json:"height"`
	Name	string	`json:"name"`
}

func NewOrderKeeper(cdb *couchdb.GoCouchDB, key sdk.StoreKey) OrderKeeper {
	return OrderKeeper{
		cdb:		cdb,
		StoreKey:	key,
	}
}

func (ok *OrderKeeper) WaitForReady(ctx sdk.Context) {
	for {
		//更新store
		//ctx.LoadLatestVersion()
		orderbook, err := ok.GetOrderBook(ctx)
		if err != nil {
			panic(err)
		}

		if ok.isReady(orderbook, ctx.ChainID(), ctx.BlockHeight()) {
			ok.UpdateOrderBook(ctx, orderbook, nil)
			return
		}
		time.Sleep(SleepTime)
		if err != nil {
			panic(err)
		}
	}
}

func (ok *OrderKeeper) UpdateOrderBook(ctx sdk.Context, orderbook OrderBook, actions *Actions) {
	if actions != nil {
		orderbook.Actions = append(orderbook.Actions, *actions)
	}

	for i := 0; i < len(orderbook.Lists); i++ {
		if orderbook.Lists[i].Name == ctx.ChainID(){
			orderbook.Lists[i].Height = ctx.BlockHeight()
			orderbook.Current.Index = i
			orderbook.Current.State = StateDone
			break
		}
	}
	//handler actions
	if orderbook.Current.Index == 0 && orderbook.Actions != nil {
		var actions []Actions
		for k, v := range orderbook.Actions {
			if v.Type == "ADD" && ctx.BlockHeight() == v.Height {
				list := Lists{
					Name:   v.Name,
					Height: 0,
				}
				orderbook.Lists = append(orderbook.Lists, list)
				length := len(orderbook.Actions)
				if length - 1 > k {
					orderbook.Actions = orderbook.Actions[k+1:]
				} else {
					orderbook.Actions = actions
				}
			}
		}
	}
	ok.SetOrderBook(ctx, orderbook)
	return
}

func (ok *OrderKeeper) GetOrderBook(ctx sdk.Context) (OrderBook, error) {
	store := ctx.KVStore(ok.StoreKey)
	var orderbook OrderBook
	bz := store.Get([]byte(OrderBookKey))
	err := ModuleCdc.UnmarshalJSON(bz, &orderbook)
	return orderbook, err
}

func (ok *OrderKeeper) ExistOrderBook(ctx sdk.Context) bool  {
	store := ctx.KVStore(ok.StoreKey)
	bz := store.Get([]byte(OrderBookKey))
	if len(bz) > 0 {
		return true
	}
	return false
}

func (ok *OrderKeeper) SetOrderBook(ctx sdk.Context, orderbook OrderBook)  {
	store := ctx.KVStore(ok.StoreKey)
	bz, err := ModuleCdc.MarshalJSON(orderbook)
	if err != nil {
		panic(err)
	}
	store.Set([]byte(OrderBookKey), bz)
}

func (ok *OrderKeeper) isReady(orderbook OrderBook, shardID string, height int64) bool {
	if orderbook.Current.State == StateInit {
		if orderbook.Lists[0].Name == shardID {
			return true
		} else {
			return false
		}
	}
	var nextIndex int
	if orderbook.Current.Index == len(orderbook.Lists) - 1 {
		nextIndex = 0
	} else {
		nextIndex = orderbook.Current.Index + 1
	}
	if orderbook.Lists[nextIndex].Height + 1 == height &&
		orderbook.Current.State == StateDone &&
		orderbook.Lists[nextIndex].Name == shardID {
		return true
	}else {
		return false
	}
}