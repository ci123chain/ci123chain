package redis

import (
	"encoding/hex"
	"github.com/ci123chain/ci123chain/pkg/libs"
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
	_, err := libs.RetryI(0, func(retryTimes int) (bytes interface{}, e error) {
		_, err := rb.batch.Commit()
		if err != nil {
			rb.rdb.lg.Error("batch write failed", "Method", "Write",
				"Retry times", retryTimes, "error", err.Error())
			return nil, err
		}else {
			return nil, nil
		}
	})
	return err
}


func (rb *redisBatch) WriteSync() error {
	return rb.Write()
}


func (rb *redisBatch) Close() error {
	rb = nil
	return nil
}