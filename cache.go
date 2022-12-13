package genericcache

import (
	"fmt"
	"sync"
	"time"
)

var ErrItemNotInCache error = fmt.Errorf("item not in cache")

type CacheValue interface {
	GetExpiryDate() int64
}

type GenericCache[C CacheValue] struct {
	CacheMap map[string]C
	stop     chan struct{}
	wg       *sync.WaitGroup
	mu       sync.RWMutex
}

func NewGenericCache[C CacheValue](cacheInterval time.Duration) *GenericCache[C] {
	cache := &GenericCache[C]{
		CacheMap: map[string]C{},
		stop:     make(chan struct{}),
		wg:       new(sync.WaitGroup),
	}

	cache.wg.Add(1)
	go func() {
		defer cache.wg.Done()
		cache.cleanUpLoop(cacheInterval)
	}()

	return cache
}

func (c *GenericCache[C]) Get(key string) (C, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.CacheMap[key]
	if !ok {
		return value, ErrItemNotInCache
	}

	return value, nil
}

func (c *GenericCache[C]) Set(key string, value C) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CacheMap[key] = value

	return nil
}

func (c *GenericCache[C]) cleanUpLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-t.C:
			c.mu.Lock()
			for key, value := range c.CacheMap {
				if value.GetExpiryDate() <= time.Now().Unix() {
					delete(c.CacheMap, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
