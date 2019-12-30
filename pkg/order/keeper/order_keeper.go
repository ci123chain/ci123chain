package keeper

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"github.com/tanhuiya/ci123chain/pkg/params/subspace"
	"strconv"
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
	Turns 	[]Turns 	`json:"turns"`

	Current	Current 	`json:"current"`

	Events	[]Events 	`json:"events"`
}

type Turns struct {
	Name 	string 	`json:"name"`
	Height	int64	`json:"height"`
}

type Current struct {
	Index	int		`json:"index"`
	State	string	`json:"state"`
}

type Events struct {
	Type	string	`json:"type"`
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
		if ok.isReady(orderbook, shardID, height) {
			err := ok.UpdateOrderBook(orderbook, rev, shardID, height, StateProcessing) //change state to processing
			if err != nil {  // other peer is processing, wait
				er := err.(*couchdb.Error)
				if er.Reason == "Document update conflict." {
					ok.waitOtherPeer(orderbook, shardID, height)
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

func (ok *OrderKeeper) UpdateOrderBook(orderBook OrderBook, rev, shardID string, height int64, state string) error {

	for i := 0; i < len(orderBook.Turns); i++ {
		if orderBook.Turns[i].Name == shardID{
			orderBook.Turns[i].Height = height
			orderBook.Current.Index = i
			orderBook.Current.State = state
			break
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
	if rev == "" || rev[0:1] == "1-"{
		ok.cdb.Set([]byte(OrderBookKey), orderBytes)
	}
}

func (ok *OrderKeeper) isReady(orderbook OrderBook, shardID string, height int64) bool {
	ID, err := strconv.Atoi(shardID[5:])
	if err != nil {
		panic(err)
	}
	if orderbook.Current.State == StateInit {
		if ID == 1 {
			return true
		} else {
			return false
		}
	}
	var nextIndex int
	if orderbook.Current.Index == len(orderbook.Turns) - 1 {
		nextIndex = 0
	} else {
		nextIndex = orderbook.Current.Index + 1
	}
	if orderbook.Turns[nextIndex].Height + 1 == height &&
		orderbook.Current.State == StateDone &&
		orderbook.Turns[nextIndex].Name == shardID {
		return true
	}else {
		return false
	}
}

func (ok *OrderKeeper) waitOtherPeer(orderBook OrderBook, shardID string, height int64) {
	for {
		_, orderBook1 := ok.GetOrderBook()
		if orderBook1.Turns[orderBook.Current.Index].Name == shardID &&
			orderBook1.Turns[orderBook.Current.Index].Height == height &&
			orderBook1.Current.State == StateProcessing {
			time.Sleep(SleepTime)
			continue
		}
		ok.IsDeal = false
		return
	}
}