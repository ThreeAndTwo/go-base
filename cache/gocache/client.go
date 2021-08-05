package gocache

import (
	"github.com/deng00/go-base/cache"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

var _ cache.Interface = &Client{}

type Client struct {
	cache *gocache.Cache
}

func (c *Client) Exist(key string) (bool, error) {
	_, ok := c.cache.Get(key)
	return ok, nil
}

func (c *Client) Get(key string) (interface{}, error) {
	val, ok := c.cache.Get(key)
	if !ok {
		return nil, cache.ErrKeyNotFound
	}
	return val, nil
}

func (c *Client) Set(key string, value interface{}) error {
	c.cache.SetDefault(key, value)
	return nil
}
func (c *Client) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	c.cache.Set(key, value, expiration)
	return nil
}
func New(defaultExpiration time.Duration) *Client {
	c := gocache.New(defaultExpiration, defaultExpiration)
	return &Client{cache: c}
}
