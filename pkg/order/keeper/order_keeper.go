package keeper

import (
	"encoding/json"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/app/types"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"github.com/tanhuiya/ci123chain/pkg/logger"
	"github.com/tanhuiya/ci123chain/pkg/params/subspace"
	"time"
)

const SleepTime = 1 * time.Second
const StateProcessing = "Processing"
const StateDone = "Done"
const StateInit = "Init"
const OrderBookKey = "OrderBook"
type OrderKeeper struct {
	cdb 		*couchdb.GoCouchDB
	paramSubspace subspace.Subspace
	IsDeal		bool
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

func NewOrderKeeper(cdb *couchdb.GoCouchDB) OrderKeeper {
	return OrderKeeper{
		cdb:		cdb,
	}
}

func (ok *OrderKeeper) WaitForReady(shardID string, height int64) {
	for {
		rev, orderbook := ok.GetOrderBook()
		logger.GetLogger().Info(fmt.Sprintf("waiting for minute %v", orderbook))
		if ok.isReady(orderbook, shardID, height) {
			err := ok.UpdateOrderBook(orderbook, rev, shardID, height, StateProcessing) //change state to processing
			if err != nil {  // other peer is processing, wait
				er := err.(*couchdb.Error)
				if er.Reason == "Document update conflict." {
					ok.waitOtherPeer(shardID, height)
					return
				} else {
					panic(err)
				}
			} else { // our turn
				ok.IsDeal = true
				return
			}
		}
		time.Sleep(SleepTime)
	}
}

func (ok *OrderKeeper) SetEventBook(orderBook OrderBook) {
	orderBytes, err := json.Marshal(orderBook)
	if err != nil {
		panic(err)
	}
	ok.cdb.Set([]byte(OrderBookKey), orderBytes)
}

func (ok *OrderKeeper) UpdateOrderBook(orderBook OrderBook, rev, shardID string, height int64, state string) error {

	for i := 0; i < len(orderBook.Lists); i++ {
		if orderBook.Lists[i].Name == shardID{
			orderBook.Lists[i].Height = height
			orderBook.Current.Index = i
			orderBook.Current.State = state
			break
		}
	}
	//handler actions
	if orderBook.Current.Index == 0 && orderBook.Actions != nil {
		var actions []Actions
		for k, v := range orderBook.Actions {
			if v.Type == "ADD" && height == v.Height {
				list := Lists{
					Name:   v.Name,
					Height: 0,
				}
				orderBook.Lists = append(orderBook.Lists, list)
				length := len(orderBook.Actions)
				if length - 1 > k {
					orderBook.Actions = orderBook.Actions[k+1:]
				} else {
					orderBook.Actions = actions
				}
			}
		}


	}

	obBytes, err := json.Marshal(orderBook)
	if err != nil {
		return err
	}
	_, err = ok.cdb.SetRev([]byte(OrderBookKey), obBytes, rev)
	if err != nil {
		return err
	}
	return nil
}

func (ok *OrderKeeper) GetOrderBook() (string, OrderBook) {
	var ob OrderBook
	rev, obBytes := ok.cdb.GetRevAndValue([]byte(OrderBookKey))
	err := json.Unmarshal(obBytes, &ob)
	if err != nil {
		panic(err)
	}
	return rev, ob
}

func (ok *OrderKeeper) SetOrderBook(orderBook OrderBook) {
	orderBytes, err := json.Marshal(orderBook)
	if err != nil {
		panic(err)
	}
	rev := ok.cdb.GetRev([]byte(OrderBookKey))
	if rev == "" {
		ok.cdb.Set([]byte(OrderBookKey), orderBytes)
	}
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

func (ok *OrderKeeper) waitOtherPeer(shardID string, height int64) {
	for {
		key := fmt.Sprintf(types.CommitInfoKeyFmt, height)
		commitID := ok.cdb.Get([]byte(key))
		if commitID != nil {
			ok.IsDeal = false
			return
		}
	}
}