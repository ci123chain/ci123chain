package keeper

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/couchdb"
	"strconv"
	"time"
)

const SleepTime = 1 * time.Second
const StateProcessing = "Processing"
const StateDone = "Done"

type OrderKeeper struct {
	cdb 		*couchdb.GoCouchDB
	IsDeal		bool
}

type OrderBook struct {
	HeightNow	int64
	ShardNow	string
	TotalShards	int
	State		string
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
			err := ok.UpdateOrderBook(orderbook.TotalShards, rev, shardID, height, StateProcessing) //change state to processing
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

func (ok *OrderKeeper) UpdateOrderBook(totalShards int, rev, shardID string, height int64, state string) error {
	orderBook := OrderBook{
		HeightNow:   height,
		ShardNow:    shardID,
		TotalShards: totalShards,
		State:       state,
	}
	obBytes, err := json.Marshal(orderBook)
	if err != nil {
		return err
	}
	_, err = ok.cdb.SetRev([]byte("OrderBook"), obBytes, rev)
	if err != nil {
		return err
	}
	return nil
}

func (ok *OrderKeeper) GetOrderBook() (string, OrderBook) {
	var ob OrderBook
	rev, obBytes := ok.cdb.GetRevAndValue([]byte("OrderBook"))
	err := json.Unmarshal(obBytes, &ob)
	if err != nil {
		panic(err)
	}
	return rev, ob
}

//
func (ok *OrderKeeper) isReady(orderbook OrderBook, shardID string, height int64) bool {
	var preID int
	var heightNow int64
	ID, err := strconv.Atoi(shardID[5:])
	if err != nil {
		panic(err)
	}
	if orderbook.HeightNow == 0 {
		if ID == 1 {
			return true
		} else {
			return false
		}
	}
	if ID == 1 {
		preID = orderbook.TotalShards
		heightNow = height - 1
	} else {
		preID = ID - 1
		heightNow = height
	}
	preShardID := "Shard" + strconv.Itoa(preID)

	if orderbook.ShardNow == preShardID && orderbook.HeightNow == heightNow && orderbook.State == StateDone {
		return true
	}else{
		return false
	}
}

func (ok *OrderKeeper) waitOtherPeer(orderBook OrderBook, shardID string, height int64) {
	for {
		_, orderBook1 := ok.GetOrderBook()
		if orderBook1.ShardNow == shardID && orderBook1.HeightNow == height && orderBook1.State == StateProcessing {
			time.Sleep(SleepTime)
			continue
		}
		ok.IsDeal = false
		return
	}
}