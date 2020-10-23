package keeper

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/couchdb"
	"github.com/ci123chain/ci123chain/pkg/order/types"
	"github.com/ci123chain/ci123chain/pkg/params/subspace"
	"time"
)

type OrderKeeper struct {
	Cdb 		*couchdb.GoCouchDB
	StoreKey	sdk.StoreKey
	paramSubspace subspace.Subspace
	AccountKeeper  account.AccountKeeper
}

func NewOrderKeeper(cdb *couchdb.GoCouchDB, key sdk.StoreKey, ak account.AccountKeeper) OrderKeeper {
	return OrderKeeper{
		Cdb:		cdb,
		StoreKey:	key,
		AccountKeeper:ak,
	}
}

func (ok *OrderKeeper) WaitForReady(ctx sdk.Context) {
	for {
		orderbook, err := ok.GetOrderBook(ctx)
		if err != nil {
			if err.Error() != types.NoOrderBookErr {
				panic(err)
			} else {
				time.Sleep(types.SleepTime)
				continue
			}
		}
		if ok.isReady(orderbook, ctx.ChainID(), ctx.BlockHeight()) {
			ok.UpdateOrderBook(ctx, orderbook, nil)
			return
		}
		time.Sleep(types.SleepTime)
		if err != nil {
			panic(err)
		}
	}
}

func (ok *OrderKeeper) UpdateOrderBook(ctx sdk.Context, orderbook types.OrderBook, actions *types.Actions) {
	if actions != nil {
		name := actions.Name
		for _,v := range orderbook.Lists{
			if v.Name == name{
				return
			}
		}
		for _,v := range orderbook.Actions{
			if v.Name == name{
				return
			}
		}
		orderbook.Actions = append(orderbook.Actions, *actions)
	}

	for i := 0; i < len(orderbook.Lists); i++ {
		if orderbook.Lists[i].Name == ctx.ChainID(){
			orderbook.Lists[i].Height = ctx.BlockHeight()
			orderbook.Current.Index = i
			orderbook.Current.State = types.StateCommitting
			break
		}
	}

	//handler actions
	var deleteIndex []int
	if orderbook.Current.Index == 0 && orderbook.Actions != nil {
		for k, v := range orderbook.Actions {
			if v.Type == types.OpADD && ctx.BlockHeight() == v.Height {
				list := types.Lists{
					Name:   v.Name,
					Height: 0,
				}
				orderbook.Lists = append(orderbook.Lists, list)
				deleteIndex = append(deleteIndex, k)
			}
		}
	}

	for k, v := range deleteIndex{
		length := len(orderbook.Actions)
		if length - 1 > 0 {
			orderbook.Actions = append(orderbook.Actions[:v-k],orderbook.Actions[v-k+1:]...)
		} else {
			orderbook.Actions = nil
		}
	}
	ok.SetOrderBook(ctx, orderbook)
	return
}

func (ok *OrderKeeper) GetOrderBook(ctx sdk.Context) (types.OrderBook, error) {
	store := ctx.KVStore(ok.StoreKey).Latest([]string{types.OrderBookKey})
	var orderbook types.OrderBook
	isExist := ok.ExistOrderBook(ctx)
	if !isExist {
		return orderbook, errors.New(types.NoOrderBookErr)
	}
	bz := store.Get([]byte(types.OrderBookKey))
	err := types.ModuleCdc.UnmarshalJSON(bz, &orderbook)
	return orderbook, err
}

func (ok *OrderKeeper) ExistOrderBook(ctx sdk.Context) bool  {
	store := ctx.KVStore(ok.StoreKey).Latest([]string{types.OrderBookKey})
	bz := store.Get([]byte(types.OrderBookKey))
	if len(bz) > 0 {
		return true
	}
	return false
}

func (ok *OrderKeeper) SetOrderBook(ctx sdk.Context, orderbook types.OrderBook)  {
	store := ctx.KVStore(ok.StoreKey)
	bz, err := types.ModuleCdc.MarshalJSON(orderbook)
	if err != nil {
		panic(err)
	}
	store.Set([]byte(types.OrderBookKey), bz)
}

func (ok *OrderKeeper) isReady(orderbook types.OrderBook, shardID string, height int64) bool {
	if orderbook.Current.State == types.StateInit {
		if orderbook.Lists[0].Name == shardID {
			return true
		} else {
			return false
		}
	}

	//handle crash
	if orderbook.Lists[orderbook.Current.Index].Name == shardID &&
		orderbook.Current.State == types.StateCommitting {
		if orderbook.Lists[orderbook.Current.Index].Height + 1 == height {
			orderbook.Current.State = types.StateDone
			orderBytes, _ :=types.ModuleCdc.MarshalJSON(orderbook)
			ok.Cdb.Set(sdk.NewPrefixedKey([]byte(types.StoreKey), []byte(types.OrderBookKey)), orderBytes)
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
		orderbook.Current.State == types.StateDone &&
		orderbook.Lists[nextIndex].Name == shardID {
		return true
	}else {
		return false
	}
}