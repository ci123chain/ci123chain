package redis

import (
	"encoding/binary"
	"errors"
	"github.com/go-redis/redis/v8"
)

type RaftRedisClient struct {
	*redis.Client
}

func NewRaftRedisClient(opt *redis.Options) *RaftRedisClient {
	c := &RaftRedisClient{Client: redis.NewClient(opt)}
	return c
}

// 批量提交
func (client *RaftRedisClient) NewBatch() *Batch {
	return &Batch{
		rdb: client.Client,
	}
}

// 获取前缀为 prefix 的所有键值对
// desc 为降序返回
func (client *RaftRedisClient) GetKeys(prefix string, desc bool) (KVPairs, error) {
	cmder := redis.NewStringSliceCmd(ctx, "keys", prefix+"*", "withvalues")
	if err := client.Process(ctx, cmder); err != nil {
		return nil, err
	}
	if err := cmder.Err(); err != nil {
		return nil, err
	}
	var pairs KVPairs
	result := cmder.Val()
	for i := 0; i < len(result); i += 2 {
		pairs = append(pairs, KVPair{
			Key:   result[i],
			Value: result[i+1],
		})
	}
	if desc {
		for i := 0; i < len(pairs)/2; i++ {
			pairs[i], pairs[len(pairs)-1-i] = pairs[len(pairs)-1-i], pairs[i]
		}
	}
	return pairs, nil
}

type Batch struct {
	rdb  *redis.Client
	docs []bulkDoc
}

type bulkDoc struct {
	isDelete bool
	args     []interface{}
}

func (bt *Batch) Set(key string, value interface{}) {
	bt.docs = append(bt.docs, bulkDoc{
		isDelete: false,
		args:     []interface{}{key, value},
	})
}

func (bt *Batch) Delete(key string) {
	bt.docs = append(bt.docs, bulkDoc{
		isDelete: true,
		args:     []interface{}{key},
	})
}

func (bt *Batch) Commit() (string, error) {
	if len(bt.docs) == 0 {
		return "", errors.New("empty batch")
	}
	args := []interface{}{"batch"}
	for _, doc := range bt.docs {
		if doc.isDelete {
			args = append(args, "d")
		} else {
			args = append(args, "s")
		}
		args = append(args, doc.args...)
	}
	cmder := redis.NewStringCmd(ctx, args...)
	if err := bt.rdb.Process(ctx, cmder); err != nil {
		return "", err
	}
	return cmder.Result()
}

func (bt *Batch) Discard() {
	bt.docs = nil
}

type KVPair struct {
	Key   string
	Value string
}

type KVPairs []KVPair

func (kv KVPairs) Len() int {
	return len(kv)
}

func (kv KVPairs) Less(i, j int) bool {
	KI := make([]byte, 8)
	KJ := make([]byte, 8)

	copy(KI, []byte(kv[i].Key)[1:9])
	copy(KJ, []byte(kv[j].Key)[1:9])

	powerI := binary.BigEndian.Uint64(KI)
	powerJ := binary.BigEndian.Uint64(KJ)

	return powerI > powerJ
}

func (kv KVPairs) Swap(i, j int) {
	kv[i], kv[j] = kv[j], kv[i]
}

func IsKeyNotExist(err error) bool {
	return err == redis.Nil
}
