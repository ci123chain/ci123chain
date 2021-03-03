package redis

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/go-redis/redis/v8"
	db "github.com/tendermint/tm-db"
)

var ctx = context.Background()

type RedisDB struct {
	DB    *RaftRedisClient
	lg    logger.Logger
}

func (rdb *RedisDB) ReverseIterator(start, end []byte) db.Iterator {
	return rdb.NewRedisIterator(start, true)
}

func NewRedisDB(opt *redis.Options) *RedisDB {
	return &RedisDB{DB:NewRaftRedisClient(opt)}
}

//check DB is connected
func DBIsValid(rdb *RedisDB) error {
	return rdb.DB.Ping(ctx).Err()
}


///implement DB
func (rdb *RedisDB) Get(key []byte) []byte {

	retry := 0
	for {
		if key == nil {
			return nil
		}else {
			value, err := rdb.DB.Get(ctx, hex.EncodeToString(key)).Result()
			if err != nil {
				if IsKeyNotExist(err) {
					return nil
				}else {
					rdb.lg.Info("***************Retry******************")
					rdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
					rdb.lg.Info(fmt.Sprintf("Method: Get, key: %s, id: %s", string(key), hex.EncodeToString(key)))
					rdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
					retry ++
					continue
				}
			}
			res, _ := hex.DecodeString(value)
			return res
		}
	}
}

func (rdb *RedisDB) Has(key []byte) bool {
	return rdb.Get(key) != nil
}

func (rdb *RedisDB) Set(key, value []byte) {
	retry := 0

	for {
		if key == nil {
			rdb.lg.Debug("the key which you set is empty")
			panic(errors.New(fmt.Sprintf("the key: %s , which you set is empty", hex.EncodeToString(key))))
		}
		if value == nil {
			rdb.lg.Debug("the value is empty, where you set the key", "key", hex.EncodeToString(key))
		}

		_, err :=rdb.DB.Set(ctx, hex.EncodeToString(key), hex.EncodeToString(value), 0).Result()
		if err != nil {
			rdb.lg.Info("***************Retry******************")
			rdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
			rdb.lg.Info(fmt.Sprintf("Method: Set, key: %s, id: %s", string(key), hex.EncodeToString(key)))
			rdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
			retry ++
			continue
		}else {
			return
		}
	}
}


func (rdb *RedisDB) SetSync(key, value []byte) {
	rdb.Set(key, value)
}


func (rdb *RedisDB) Delete(key []byte){
	retry := 0
	for {
		if !rdb.Has(key) {
			return
		}else {
			n, err := rdb.DB.Del(ctx, hex.EncodeToString(key)).Result()
			if err != nil {
				rdb.lg.Info("delete key failed", "key", hex.EncodeToString(key))
				rdb.lg.Info("***************Retry******************")
				rdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
				rdb.lg.Info(fmt.Sprintf("Method: Delete, key: %s, id: %s", string(key), hex.EncodeToString(key)))
				rdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
				retry++
			}else if n != 1 {
				rdb.lg.Info(fmt.Sprintf("unexpected return value, expect 1, got %d", n))
				rdb.lg.Info("***************Retry******************")
				rdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
				rdb.lg.Info(fmt.Sprintf("Method: Delete, key: %s, id: %s", string(key), hex.EncodeToString(key)))
				rdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
				retry++
			}else {
				return
			}
		}
	}
}

func (rdb *RedisDB) DeleteSync(key []byte){
	rdb.Delete(key)
}


func (rdb *RedisDB) Close() {}


func (rdb *RedisDB) NewBatch() db.Batch {
	batch := rdb.DB.NewBatch()
	return &redisBatch{rdb, batch}
}

func (rdb *RedisDB) Print() {}

func (rdb *RedisDB) Stats() map[string]string {
	return nil
}

func (rdb *RedisDB) Iterator(start, end []byte) db.Iterator {
	return rdb.NewRedisIterator(start, false)
}

type RedisIterator struct {
	rdb   *RedisDB
	results []KVPair
	cursor		int
	start		[]byte
	end			[]byte
	isReverse	bool
	valid       bool
}

func (rdb *RedisDB) NewRedisIterator(start []byte, isReserve bool) db.Iterator {
	retry := 0

	for {
		results, err := rdb.DB.GetKeys(hex.EncodeToString(start), isReserve)
		if err != nil {
			rdb.lg.Info("***************Retry******************")
			rdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
			rdb.lg.Info(fmt.Sprintf("Method: GetKeys, start: %s", hex.EncodeToString(start)))
			rdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
			retry++
		}else {
			return &RedisIterator{
				rdb:       rdb,
				results:   results,
				cursor:    0,
				start:     start,
				end:       nil,
				isReverse: isReserve,
				valid:     true,
			}
		}
	}
}

func (ri *RedisIterator) Domain() (start, end []byte) {
	return ri.start,ri.end
}

func (ri *RedisIterator) Valid() bool {

	return len(ri.results) > 0 && ri.valid
}

func (ri *RedisIterator) Next() {
	ri.assertValid()
	if ri.cursor < len(ri.results) - 1 {
		ri.cursor++
	}else {
		ri.valid = false
	}
}

func (ri *RedisIterator) Key() (key []byte) {
	ri.assertValid()

	if ri.isReverse {
		value, err := hex.DecodeString(ri.results[len(ri.results) - ri.cursor -1].Key)
		if err != nil {
			//TODO panic is OK?
			ri.rdb.lg.Info("iterator got value failed", "Error", err.Error())
			panic(err)
		}
		return value
	}

	value, err := hex.DecodeString(ri.results[ri.cursor].Key)
	if err != nil {
		//TODO panic is OK?
		ri.rdb.lg.Info("iterator got value failed", "Error", err.Error())
		panic(err)
	}
	return value
}

func (ri *RedisIterator) Value() (value []byte) {
	ri.assertValid()

	if ri.isReverse {
		value, err := hex.DecodeString(ri.results[len(ri.results) - ri.cursor -1].Value)
		if err != nil {
			//TODO panic is OK?
			ri.rdb.lg.Info("iterator got value failed", "Error", err.Error())
			panic(err)
		}
		return value
	}

	value, err := hex.DecodeString(ri.results[ri.cursor].Value)
	if err != nil {
		//TODO panic is OK?
		ri.rdb.lg.Info("iterator got value failed", "Error", err.Error())
		panic(err)
	}
	return value
}

func (ri *RedisIterator) Close() {
	ri = nil
}

func (ri *RedisIterator) assertValid() {
	if !ri.Valid() {
		panic("redisIterator is invalid")
	}
}