package redis

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/libs"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/go-redis/redis/v8"
	db "github.com/tendermint/tm-db"
)

var ctx = context.Background()

type RedisDB struct {
	DB    *RaftRedisClient
	lg    logger.Logger
}

func (rdb *RedisDB) ReverseIterator(start, end []byte) (db.Iterator, error) {
	return rdb.NewRedisIterator(start, end, true), nil
}

func NewRedisDB(opt *redis.Options) *RedisDB {
	return &RedisDB{DB:NewRaftRedisClient(opt), lg: logger.GetLogger()}
}

//check DB is connected
func DBIsValid(rdb *RedisDB) error {
	return rdb.DB.Ping(ctx).Err()
}


///implement DB
func (rdb *RedisDB) Get(key []byte) ([]byte, error) {
	res, err :=  libs.RetryI(0, func(retryTimes int) (interface{}, error) {
		if key == nil {
			return nil, nil
		}else {
			value, err := rdb.DB.Get(ctx, hex.EncodeToString(key)).Result()
			if err != nil {
				if IsKeyNotExist(err) {
					return nil, nil
				}else {
					rdb.lg.Error("db get failed", "Method", "Get", "Retry times", retryTimes, "key", string(key),
						"id", hex.EncodeToString(key), "error", err.Error())
					return nil, err
				}
			}
			res, _ := hex.DecodeString(value)
			return res, nil
		}
	})
	if res == nil {
		return nil, err
	}
	return res.([]byte), err
}

func (rdb *RedisDB) Has(key []byte) (bool, error) {
	v, err := rdb.Get(key)
	if err != nil {
		return false, err
	}
	return v != nil, nil
}

func (rdb *RedisDB) Set(key, value []byte) error {
	_, err := libs.RetryI(0, func(retryTimes int) (bytes interface{}, e error) {
		if key == nil {
			rdb.lg.Debug("the key which you set is empty")
			panic(errors.New(fmt.Sprintf("the key: %s , which you set is empty", hex.EncodeToString(key))))
		}
		if value == nil {
			rdb.lg.Debug("the value is empty, where you set the key", "key", hex.EncodeToString(key))
		}

		_, err :=rdb.DB.Set(ctx, hex.EncodeToString(key), hex.EncodeToString(value), 0).Result()
		if err != nil {
			rdb.lg.Error("db set failed", "Method", "Set", "Retry times", retryTimes, "key", string(key),
				"id", hex.EncodeToString(key), "error", err.Error())
			return nil, err
		}else {
			return nil, nil
		}
	})
	return err
}


func (rdb *RedisDB) SetSync(key, value []byte) error {
	return rdb.Set(key, value)
}


func (rdb *RedisDB) Delete(key []byte) error{
	_, err := libs.RetryI(0, func(retryTimes int) (bytes interface{}, e error) {
		v, err := rdb.Has(key)
		if err != nil {
			return nil, err
		}
		if !v {
			return nil, nil
		}else {
			n, err := rdb.DB.Del(ctx, hex.EncodeToString(key)).Result()
			if err != nil {
				rdb.lg.Error("delete key failed", "Method", "Delete", "Retry times", retryTimes, "key", string(key),
					"id", hex.EncodeToString(key), "error", err.Error())
				return nil, err
			}else if n != 1 {
				rdb.lg.Error(fmt.Sprintf("unexpected return value, expect 1, got %d", n), "Method", "Delete",
					"Retry times", retryTimes, "key", string(key), "id", hex.EncodeToString(key))
				return nil, fmt.Errorf("unexpected return value, expect 1, got %d", n)
			}else {
				return nil, nil
			}
		}
	})
	return err
}

func (rdb *RedisDB) DeleteSync(key []byte) error {
	return rdb.Delete(key)
}


func (rdb *RedisDB) Close() error {
	rdb = nil
	return nil
}


func (rdb *RedisDB) NewBatch() db.Batch {
	batch := rdb.DB.NewBatch()
	return &redisBatch{rdb, batch}
}

func (rdb *RedisDB) Print() error {
	return nil
}

func (rdb *RedisDB) Stats() map[string]string {
	return nil
}

func (rdb *RedisDB) Iterator(start, end []byte) (db.Iterator, error) {
	return rdb.NewRedisIterator(start, end, false), nil
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

func (ri *RedisIterator) Error() error {
	return nil
}

func (rdb *RedisDB) NewRedisIterator(start, end []byte, isReserve bool) db.Iterator {
	iterator, _ := libs.RetryI(0, func(retryTimes int) (res interface{}, err error) {
		results, err := rdb.DB.GetKeys(hex.EncodeToString(start), isReserve)
		if err != nil {
			rdb.lg.Error("db get keys failed", "Method", "Get", "Retry times", retryTimes, "keys", string(start),
				"id", hex.EncodeToString(start), "error", err.Error())
			return nil, err
		}else {
			return &RedisIterator{
				rdb:       rdb,
				results:   results,
				cursor:    0,
				start:     start,
				end:       end,
				isReverse: isReserve,
				valid:     true,
			}, nil
		}
	})
	return iterator.(db.Iterator)
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

func (ri *RedisIterator) Close() error {
	ri = nil
	return nil
}

func (ri *RedisIterator) assertValid() {
	if !ri.Valid() {
		panic("redisIterator is invalid")
	}
}