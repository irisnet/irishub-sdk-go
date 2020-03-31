package cache

import (
	"time"

	"github.com/bluele/gcache"
)

type Cache interface {
	Set(key, value interface{}) error
	SetWithExpire(key, value interface{}, expiration time.Duration) error
	Get(key interface{}) (interface{}, error)
	Remove(key interface{}) bool
}

type LRU struct {
	cache gcache.Cache
}

func NewLRU(capacity int) Cache {
	return LRU{
		cache: gcache.New(capacity).LRU().Build(),
	}
}

func (l LRU) Set(key, value interface{}) error {
	return l.cache.Set(key, value)
}

func (l LRU) SetWithExpire(key, value interface{}, expiration time.Duration) error {
	return l.cache.SetWithExpire(key, value, expiration)
}

func (l LRU) Get(key interface{}) (interface{}, error) {
	return l.cache.Get(key)
}

func (l LRU) Remove(key interface{}) bool {
	return l.cache.Remove(key)
}
