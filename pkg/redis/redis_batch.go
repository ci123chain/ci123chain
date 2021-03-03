package redis

import (
	"encoding/hex"
	"fmt"
)

type redisBatch struct {
	rdb *RedisDB
	batch *Batch
}

func (rb *redisBatch) Set(key, value []byte) {
	rb.batch.Set(hex.EncodeToString(key), hex.EncodeToString(value))
}

//TODO add batch delete.
func (rb *redisBatch) Delete(key []byte) {
	rb.batch.Delete(hex.EncodeToString(key))
}

func (rb *redisBatch) Write() {
	if rb.batch.docs == nil {
		return
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
			return
		}
	}
}


func (rb *redisBatch) WriteSync() {
	rb.Write()
}


func (rb *redisBatch) Close() {
	rb = nil
}