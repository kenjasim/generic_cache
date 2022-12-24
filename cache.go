// package genericcache is a test implementation of a generic cache library using golang's new generics
// feature. It stores values which contain an expiration date. This date is checked regularly by the
// cleanup loop which will remove the item if it has existed too long in the cache.
package genericcache

import (
	"fmt"
	"sync"
	"time"
)

// ErrItemNotInCache is an error thrown when the key provided has no corresponding
// value in the cache.
var ErrItemNotInCache error = fmt.Errorf("item not in cache")

// CacheValue is an interface defining the types of items which can be inside the cache
// it requires only one function, GetExpiryDate which is used to get the time at which the
// value must expire
type CacheValue interface {
	GetExpiryDate() int64
}

// GenericCache is a struct value holding the necessary items needed for a cache. The most important is the cache itself
// which is a map with a string key and a value thc C generic type which must implement the cache value interface. The rest
// of the items in the cache are used to control the goroutine which deals with cache cleanup and the mutex for reading from
// the cache.
type GenericCache[C CacheValue] struct {
	CacheMap map[string]C
	stop     chan struct{}
	wg       *sync.WaitGroup
	mu       sync.RWMutex
}

// NewGenericCache generates a new cache of the specified type. The funtion also starts the cleanup loop which in a seperate
// go routine to check the cache values on the specified intervals.
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

// Get returns the value corresponding to the key specified if that key exists in the map, if it does not then
// the ErrItemNotInCache error is thrown.
func (c *GenericCache[C]) Get(key string) (C, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.CacheMap[key]
	if !ok {
		return value, ErrItemNotInCache
	}

	return value, nil
}

// Set writes the specified key value pair into the cache map
func (c *GenericCache[C]) Set(key string, value C) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CacheMap[key] = value

	return nil
}

// cleanUpLoop will check with the frequency passed into the function the cache values, if one is expired then the
// item is removed from the cache.
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
