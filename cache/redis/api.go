package redis

import (
	"github.com/deng00/go-base/cache"
	"github.com/go-redis/redis"
	"time"
)

type Interface interface {
	StrOp
	MapOp
	ListOp
	SetOp
	ZSetOp
	PubSub
	Bitmap
	GetCmdable() redis.Cmdable
}
type StrOp interface {
	cache.Interface
	GetString(key string) (value string, err error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	GetSet(key string, value interface{}) (string, error)
	StrLen(key string) (int64, error)
	Append(key string, value string) (int64, error)
	SetRange(key string, offset int64, value string) (int64, error)
	GetRange(key string, start, end int64) (string, error)
	Incr(key string) (int64, error)
	IncrBy(key string, incr int64) (int64, error)
	IncrByFloat(key string, incr float64) (float64, error)
	Decr(key string) (int64, error)
	DecrBy(key string, decr int64) (int64, error)
	MSet(pairs map[string]interface{}) (string, error)
	MSetNX(pairs map[string]interface{}) (bool, error)
	MGet(keys ...string) ([]interface{}, error)
	Del(key ...string) (int64, error)
	Expire(key string, expiration time.Duration) (bool, error)
	TTL(key string) (time.Duration, error)
	Exists(keys ...string) (int64, error)
	Keys(pattern string) (results []string, err error)
	Eval(script string, keys []string, args ...interface{}) (interface{}, error)
	Close() (err error)
}
type SetOp interface {
	SAdd(key string, member interface{}) (count int64, err error)
	SIsMember(key, member string) (res bool, err error)
	SPop(key string) (string, error)
	SRandMember(key string) (string, error)
	SRem(key string, member interface{}) (count int64, err error)
	SMove(source, destination string, member interface{}) (bool, error)
	SCard(key string) (int64, error)
	SMembers(key string) (value []string, err error)
	SScan(key string, cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error)
	SInter(keys ...string) (value []string, err error)
	SInterStore(destination string, keys ...string) (int64, error)
	SUnion(keys ...string) (value []string, err error)
	SUnionStore(destination string, keys ...string) (int64, error)
	SDiff(keys ...string) (value []string, err error)
	SDiffStore(destination string, keys ...string) (int64, error)
}

type ZSetOp interface {
	ZAdd(key string, members ...redis.Z) (int64, error)
	ZScore(key, member string) (float64, error)
	ZIncrBy(key string, increment float64, member string) (float64, error)
	ZCard(key string) (int64, error)
	ZCount(key, min, max string) (int64, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRangeByScore(key string, opt redis.ZRangeBy) ([]string, error)
	ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error)
	ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error)
	ZRank(key, member string) (int64, error)
	ZRevRank(key, member string) (int64, error)
	ZRem(key string, members ...interface{}) (int64, error)
	ZRemRangeByLex(key, min, max string) (int64, error)
	ZRemRangeByRank(key string, start, stop int64) (int64, error)
	ZRemRangeByScore(key, min, max string) (int64, error)
	ZRangeByLex(key string, opt redis.ZRangeBy) ([]string, error)
	ZLexCount(key, min, max string) (int64, error)
	ZScan(key string, cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error)
	ZUnionStore(dest string, store redis.ZStore, keys ...string) (int64, error)
	ZInterStore(dest string, store redis.ZStore, keys ...string) (int64, error)
}
type MapOp interface {
	HSet(key string, field string, value interface{}) (bool, error)
	HSetNX(key string, field string, value interface{}) (bool, error)
	HGet(key string, field string) (value string, err error)
	HExists(key string, field string) (bool, error)
	HDel(key string, field ...string) (count int64, err error)
	HLen(key string) (int64, error)
	HIncrBy(key string, field string, incr int64) (int64, error)
	HIncrByFloat(key string, field string, incr float64) (float64, error)
	HMSet(key string, fields map[string]interface{}) (string, error)
	HMGet(key string, fields ...string) ([]interface{}, error)
	HKeys(key string) ([]string, error)
	HVals(key string) ([]string, error)
	HGetAll(key string) (value map[string]string, err error)
	HScan(key string, cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error)
}
type ListOp interface {
	LPush(key string, values ...interface{}) (int64, error)
	LPushX(key string, value interface{}) (int64, error)
	RPush(key string, values ...interface{}) (int64, error)
	RPushX(key string, value interface{}) (int64, error)
	LPop(key string) (string, error)
	RPop(key string) (string, error)
	RPopLPush(source, destination string) (string, error)
	LRem(key string, count int64, value interface{}) (int64, error)
	LLen(key string) (int64, error)
	LIndex(key string, index int64) (string, error)
	LInsert(key, op string, pivot, value interface{}) (int64, error)
	LSet(key string, index int64, value interface{}) (string, error)
	LRange(key string, start, stop int64) ([]string, error)
	LTrim(key string, start, stop int64) (string, error)
	BLPop(timeout time.Duration, keys ...string) ([]string, error)
	BRPop(timeout time.Duration, keys ...string) ([]string, error)
	BRPopLPush(source, destination string, timeout time.Duration) (string, error)
}
type Event struct {
	Channel string
	Payload string
}
type Subscriber func(event Event)
type PubSub interface {
	Publish(event Event) (results int64, err error)
	Subscribe(subscriber Subscriber, channel ...string)
}

type Bitmap interface {
	GetBit(key string, offset int64) (int64, error)
	SetBit(key string, offset int64, bit int) error
}
