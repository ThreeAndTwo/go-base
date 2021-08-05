package cache

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Interface interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Exist(key string) (bool, error)
	SetWithExpiration(key string, value interface{}, expiration time.Duration) error
}

// Stringify serialization any type
func Stringify(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	default:
		res, _ := json.Marshal(val)
		return string(res)
	}
}
