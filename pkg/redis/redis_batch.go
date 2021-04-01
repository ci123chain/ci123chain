package redis

import (
	"encoding/hex"
	"fmt"
)

type redisBatch struct {
	rdb *RedisDB
	batch *Batch
}

func (rb *redisBatch) Set(key, value []byte) error {
	rb.batch.Set(hex.EncodeToString(key), hex.EncodeToString(value))
	return nil
}

//TODO add batch delete.
func (rb *redisBatch) Delete(key []byte) error {
	rb.batch.Delete(hex.EncodeToString(key))
	return nil
}

func (rb *redisBatch) Write() error {
	if rb.batch.docs == nil {
		return nil
	}
	retry := 0
	for {
		_, err := rb.batch.Commit()
		if err != nil {
			rb.rdb.lg.Info("batch write failed")
			rb.rdb.lg.Info("***************Retry******************")
			rb.rdb.lg.Info(fmt.Sprintf("Retry: %d", retry))
			rb.rdb.lg.Info(fmt.Sprintf("Method: Write"))
			rb.rdb.lg.Error(fmt.Sprintf("Error: %s", err.Error()))
			retry++
		}else {
			return nil
		}
	}
}


func (rb *redisBatch) WriteSync() error {
	return rb.Write()
}


func (rb *redisBatch) Close() error {
	rb = nil
	return nil
}