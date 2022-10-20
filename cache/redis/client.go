package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/deng00/go-base/cache"
	"github.com/go-redis/redis/v8"
	"time"
)

// Config redis client config
type Config struct {
	Addr     string `mapstructure:"addr"`
	Pass     string `mapstructure:"pass"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func (c Config) Check() error {
	if c.Addr == "" {
		return fmt.Errorf("invalid addr")
	}
	return nil
}

var _ cache.Interface = &Redis{}
var _ Interface = &Redis{}

// Redis redis client
type Redis struct {
	config *Config
	client *redis.Client
	ctx    context.Context
	exitCh chan struct{}
}

func (r *Redis) GetCmdable() redis.Cmdable {
	return r.client
}

func (r *Redis) Exist(key string) (bool, error) {
	res, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	if res > 0 {
		return true, nil
	}
	return false, nil
}

func (r *Redis) GetString(key string) (value string, err error) {
	result, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", cache.ErrKeyNotFound
	}
	return result, err
}

func (r *Redis) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(r.ctx, key, cache.Stringify(value), expiration).Result()
}
func (r *Redis) GetSet(key string, value interface{}) (string, error) {
	return r.client.GetSet(r.ctx, key, cache.Stringify(value)).Result()
}
func (r *Redis) StrLen(key string) (int64, error) {
	return r.client.StrLen(r.ctx, key).Result()
}
func (r *Redis) Append(key string, value string) (int64, error) {
	return r.client.Append(r.ctx, key, value).Result()
}
func (r *Redis) SetRange(key string, offset int64, value string) (int64, error) {
	return r.client.SetRange(r.ctx, key, offset, value).Result()
}
func (r *Redis) GetRange(key string, start, end int64) (string, error) {
	return r.client.GetRange(r.ctx, key, start, end).Result()
}
func (r *Redis) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}
func (r *Redis) IncrBy(key string, incr int64) (int64, error) {
	return r.client.IncrBy(r.ctx, key, incr).Result()
}
func (r *Redis) IncrByFloat(key string, incr float64) (float64, error) {
	return r.client.IncrByFloat(r.ctx, key, incr).Result()
}
func (r *Redis) Decr(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}
func (r *Redis) DecrBy(key string, decr int64) (int64, error) {
	return r.client.DecrBy(r.ctx, key, decr).Result()
}
func (r *Redis) MSet(pairs map[string]interface{}) (string, error) {
	var strPairs []interface{}
	for k, v := range pairs {
		strPairs = append(strPairs, k)
		strPairs = append(strPairs, cache.Stringify(v))
	}
	return r.client.MSet(r.ctx, strPairs...).Result()
}
func (r *Redis) MSetNX(pairs map[string]interface{}) (bool, error) {
	var strPairs []interface{}
	for k, v := range pairs {
		strPairs = append(strPairs, k)
		strPairs = append(strPairs, cache.Stringify(v))
	}
	return r.client.MSetNX(r.ctx, strPairs...).Result()
}
func (r *Redis) MGet(keys ...string) ([]interface{}, error) {
	return r.client.MGet(r.ctx, keys...).Result()
}
func (r *Redis) Del(key ...string) (int64, error) {
	return r.client.Del(r.ctx, key...).Result()
}
func (r *Redis) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(r.ctx, script, keys, args...).Result()
}
func (r *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(r.ctx, key, expiration).Result()
}

func (r *Redis) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

func (r *Redis) GetBit(key string, offset int64) (int64, error) {
	return r.client.GetBit(r.ctx, key, offset).Result()
}

func (r *Redis) SetBit(key string, offset int64, value int) error {
	return r.client.SetBit(r.ctx, key, offset, value).Err()
}

func (r *Redis) Exists(keys ...string) (int64, error) {
	return r.client.Exists(r.ctx, keys...).Result()
}

func (r *Redis) Get(key string) (interface{}, error) {
	return r.GetString(key)
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) (err error) {
	return r.client.Set(r.ctx, key, cache.Stringify(value), expiration).Err()
}
func (r *Redis) SetWithExpiration(key string, value interface{}, expiration time.Duration) (err error) {
	return r.client.Set(r.ctx, key, cache.Stringify(value), expiration).Err()
}
func (r *Redis) HSet(key string, field string, value interface{}) (int64, error) {
	return r.client.HSet(r.ctx, key, field, cache.Stringify(value)).Result()
}
func (r *Redis) HSetNX(key string, field string, value interface{}) (bool, error) {
	return r.client.HSetNX(r.ctx, key, field, cache.Stringify(value)).Result()
}
func (r *Redis) HGet(key string, field string) (value string, err error) {
	result, err := r.client.HGet(r.ctx, key, field).Result()
	if err == redis.Nil {
		return "", cache.ErrKeyNotFound
	}
	return result, err
}
func (r *Redis) HExists(key string, field string) (bool, error) {
	return r.client.HExists(r.ctx, key, field).Result()
}
func (r *Redis) HDel(key string, field ...string) (count int64, err error) {
	return r.client.HDel(r.ctx, key, field...).Result()
}
func (r *Redis) HLen(key string) (int64, error) {
	return r.client.HLen(r.ctx, key).Result()
}
func (r *Redis) HIncrBy(key string, field string, incr int64) (int64, error) {
	return r.client.HIncrBy(r.ctx, key, field, incr).Result()
}
func (r *Redis) HIncrByFloat(key string, field string, incr float64) (float64, error) {
	return r.client.HIncrByFloat(r.ctx, key, field, incr).Result()
}
func (r *Redis) HMSet(key string, fields map[string]interface{}) (bool, error) {
	strFields := make(map[string]interface{})
	for k, v := range fields {
		strFields[k] = cache.Stringify(v)
	}
	return r.client.HMSet(r.ctx, key, strFields).Result()
}
func (r *Redis) HMGet(key string, fields ...string) ([]interface{}, error) {
	return r.client.HMGet(r.ctx, key, fields...).Result()
}
func (r *Redis) HKeys(key string) ([]string, error) {
	return r.client.HKeys(r.ctx, key).Result()
}
func (r *Redis) HVals(key string) ([]string, error) {
	return r.client.HVals(r.ctx, key).Result()
}
func (r *Redis) HScan(key string, cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error) {
	return r.client.HScan(r.ctx, key, cursor, match, count).Result()
}
func (r *Redis) HGetAll(key string) (value map[string]string, err error) {
	return r.client.HGetAll(r.ctx, key).Result()
}
func (r *Redis) LPush(key string, values ...interface{}) (int64, error) {
	strValues := make([]interface{}, len(values))
	for i, v := range values {
		strValues[i] = cache.Stringify(v)
	}
	return r.client.LPush(r.ctx, key, values...).Result()
}
func (r *Redis) LPushX(key string, value interface{}) (int64, error) {
	return r.client.LPushX(r.ctx, key, cache.Stringify(value)).Result()
}
func (r *Redis) RPush(key string, values ...interface{}) (int64, error) {
	strValues := make([]interface{}, len(values))
	for i, v := range values {
		strValues[i] = cache.Stringify(v)
	}
	return r.client.RPush(r.ctx, key, values...).Result()
}
func (r *Redis) RPushX(key string, value interface{}) (int64, error) {
	return r.client.RPushX(r.ctx, key, cache.Stringify(value)).Result()
}
func (r *Redis) LPop(key string) (string, error) {
	return r.client.LPop(r.ctx, key).Result()
}
func (r *Redis) RPop(key string) (string, error) {
	return r.client.RPop(r.ctx, key).Result()
}
func (r *Redis) RPopLPush(source, destination string) (string, error) {
	return r.client.RPopLPush(r.ctx, source, destination).Result()
}
func (r *Redis) LRem(key string, count int64, value interface{}) (int64, error) {
	return r.client.LRem(r.ctx, key, count, cache.Stringify(value)).Result()
}
func (r *Redis) LLen(key string) (int64, error) {
	return r.client.LLen(r.ctx, key).Result()
}
func (r *Redis) LIndex(key string, index int64) (string, error) {
	return r.client.LIndex(r.ctx, key, index).Result()
}
func (r *Redis) LInsert(key, op string, pivot, value interface{}) (int64, error) {
	return r.client.LInsert(r.ctx, key, op, pivot, cache.Stringify(value)).Result()
}
func (r *Redis) LSet(key string, index int64, value interface{}) (string, error) {
	return r.client.LSet(r.ctx, key, index, cache.Stringify(value)).Result()
}
func (r *Redis) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, start, stop).Result()
}
func (r *Redis) LTrim(key string, start, stop int64) (string, error) {
	return r.client.LTrim(r.ctx, key, start, stop).Result()
}
func (r *Redis) BLPop(timeout time.Duration, keys ...string) ([]string, error) {
	return r.client.BLPop(r.ctx, timeout, keys...).Result()
}
func (r *Redis) BRPop(timeout time.Duration, keys ...string) ([]string, error) {
	return r.client.BRPop(r.ctx, timeout, keys...).Result()
}
func (r *Redis) BRPopLPush(source, destination string, timeout time.Duration) (string, error) {
	return r.client.BRPopLPush(r.ctx, source, destination, timeout).Result()
}
func (r *Redis) SAdd(key string, member interface{}) (count int64, err error) {
	return r.client.SAdd(r.ctx, key, cache.Stringify(member)).Result()
}

func (r *Redis) SIsMember(key, member string) (res bool, err error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}
func (r *Redis) SPop(key string) (string, error) {
	return r.client.SPop(r.ctx, key).Result()
}
func (r *Redis) SRandMember(key string) (string, error) {
	return r.client.SRandMember(r.ctx, key).Result()
}
func (r *Redis) SRem(key string, member interface{}) (count int64, err error) {
	return r.client.SRem(r.ctx, key, cache.Stringify(member)).Result()
}
func (r *Redis) SMove(source, destination string, member interface{}) (bool, error) {
	return r.client.SMove(r.ctx, source, destination, cache.Stringify(member)).Result()
}
func (r *Redis) SCard(key string) (int64, error) {
	return r.client.SCard(r.ctx, key).Result()
}
func (r *Redis) SMembers(key string) (value []string, err error) {
	return r.client.SMembers(r.ctx, key).Result()
}
func (r *Redis) SScan(key string, cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error) {
	return r.client.SScan(r.ctx, key, cursor, match, count).Result()
}
func (r *Redis) SInter(keys ...string) (value []string, err error) {
	return r.client.SInter(r.ctx, keys...).Result()
}
func (r *Redis) SInterStore(destination string, keys ...string) (int64, error) {
	return r.client.SInterStore(r.ctx, destination, keys...).Result()
}
func (r *Redis) SUnion(keys ...string) (value []string, err error) {
	return r.client.SUnion(r.ctx, keys...).Result()
}
func (r *Redis) SUnionStore(destination string, keys ...string) (int64, error) {
	return r.client.SUnionStore(r.ctx, destination, keys...).Result()
}
func (r *Redis) SDiff(keys ...string) (value []string, err error) {
	return r.client.SDiff(r.ctx, keys...).Result()
}
func (r *Redis) SDiffStore(destination string, keys ...string) (int64, error) {
	return r.client.SDiffStore(r.ctx, destination, keys...).Result()
}
func (r *Redis) ZAdd(key string, members ...*redis.Z) (int64, error) {
	return r.client.ZAdd(r.ctx, key, members...).Result()
}
func (r *Redis) ZScore(key, member string) (float64, error) {
	return r.client.ZScore(r.ctx, key, member).Result()
}
func (r *Redis) ZIncrBy(key string, increment float64, member string) (float64, error) {
	return r.client.ZIncrBy(r.ctx, key, increment, member).Result()
}
func (r *Redis) ZCard(key string) (int64, error) {
	return r.client.ZCard(r.ctx, key).Result()
}
func (r *Redis) ZCount(key, min, max string) (int64, error) {
	return r.client.ZCount(r.ctx, key, min, max).Result()
}
func (r *Redis) ZRange(key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(r.ctx, key, start, stop).Result()
}

func (r *Redis) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRangeWithScores(r.ctx, key, start, stop).Result()
}

func (r *Redis) ZRevRange(key string, start, stop int64) ([]string, error) {
	return r.client.ZRevRange(r.ctx, key, start, stop).Result()
}

func (r *Redis) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRevRangeWithScores(r.ctx, key, start, stop).Result()
}

func (r *Redis) ZRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByScore(r.ctx, key, opt).Result()
}

func (r *Redis) ZRevRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRevRangeByScore(r.ctx, key, opt).Result()
}

func (r *Redis) ZRank(key, member string) (int64, error) {
	return r.client.ZRank(r.ctx, key, member).Result()
}
func (r *Redis) ZRevRank(key, member string) (int64, error) {
	return r.client.ZRevRank(r.ctx, key, member).Result()
}
func (r *Redis) ZRem(key string, members ...interface{}) (int64, error) {
	strMembers := make([]interface{}, len(members))
	for i, v := range members {
		strMembers[i] = cache.Stringify(v)
	}
	return r.client.ZRem(r.ctx, key, strMembers...).Result()
}
func (r *Redis) ZRemRangeByLex(key, min, max string) (int64, error) {
	return r.client.ZRemRangeByLex(r.ctx, key, min, max).Result()
}
func (r *Redis) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return r.client.ZRemRangeByRank(r.ctx, key, start, stop).Result()
}
func (r *Redis) ZRemRangeByScore(key, min, max string) (int64, error) {
	return r.client.ZRemRangeByScore(r.ctx, key, min, max).Result()
}
func (r *Redis) ZRangeByLex(key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByLex(r.ctx, key, opt).Result()
}
func (r *Redis) ZLexCount(key, min, max string) (int64, error) {
	return r.client.ZLexCount(r.ctx, key, min, max).Result()
}
func (r *Redis) ZScan(key string, cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error) {
	return r.client.ZScan(r.ctx, key, cursor, match, count).Result()
}
func (r *Redis) ZUnionStore(dest string, store *redis.ZStore) (int64, error) {
	return r.client.ZUnionStore(r.ctx, dest, store).Result()
}
func (r *Redis) ZInterStore(dest string, store *redis.ZStore) (int64, error) {
	return r.client.ZInterStore(r.ctx, dest, store).Result()
}
func (r *Redis) Keys(pattern string) (results []string, err error) {
	return r.client.Keys(r.ctx, pattern).Result()
}
func (r *Redis) Publish(event Event) (results int64, err error) {
	return r.client.Publish(r.ctx, event.Channel, event.Payload).Result()
}
func (r *Redis) Subscribe(subscriber Subscriber, channel ...string) {
	pubSub := r.client.Subscribe(r.ctx, channel...)
	subChan := pubSub.Channel()
	go func() {
		for {
			select {
			case <-r.exitCh:
				_ = pubSub.Close()
				return
			case event := <-subChan:
				subscriber(Event{
					Channel: event.Channel,
					Payload: event.Payload,
				})
			}
		}
	}()
}

// New create new redis client
func New(config *Config) (client *Redis, err error) {
	ctx := context.TODO()
	if err := config.Check(); err != nil {
		return nil, fmt.Errorf("invalid config:%s", err)
	}
	cli := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Pass,
		DB:       config.DB,
		PoolSize: config.PoolSize,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	_, err = cli.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	client = &Redis{
		config: config,
		client: cli,
		ctx:    ctx,
		exitCh: make(chan struct{}),
	}
	return
}
func (r *Redis) Close() (err error) {
	close(r.exitCh)
	return r.client.Close()
}

func (r *Redis) GetConfig() *Config {
	return r.config
}
