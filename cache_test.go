package genericcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCacheValue which implements the interface for the cache value
type TestCacheValue struct {
	value  string
	expiry time.Time
}

// GetExpiryDate - Returns the expiry date for the
func (t *TestCacheValue) GetExpiryDate() int64 {
	return t.expiry.Unix()
}

func TestGetValue(t *testing.T) {

	testcases := []struct {
		name          string
		key           string
		value         *TestCacheValue
		expectedError error
		cacheDuration time.Duration
		expiredTest   bool
	}{
		{
			name: "value in the cache",
			key:  "foo",
			value: &TestCacheValue{
				value:  "foo",
				expiry: time.Now().Add(24 * time.Hour),
			},
			expectedError: nil,
			cacheDuration: 24 * time.Hour,
		},
		{
			name:          "no value in the cache",
			key:           "foo",
			value:         nil,
			expectedError: ItemNotInCache,
			cacheDuration: 24 * time.Hour,
		},
		{
			name: "value in the cache but expired",
			key:  "foo",
			value: &TestCacheValue{
				value:  "foo",
				expiry: time.Now().Add(1 * time.Microsecond),
			},
			expectedError: ItemNotInCache,
			cacheDuration: 1 * time.Second,
			expiredTest:   true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {

			// Create new generic cache for the test value with a 24 hour cleanup duration (default)
			cache := NewGenericCache[*TestCacheValue](testcase.cacheDuration)
			if testcase.value != nil {
				cache.CacheMap[testcase.key] = testcase.value
			}

			if testcase.expiredTest {
				// Wait and set the value to nil as we are testing that the item is removed
				time.Sleep(1 * time.Second)
				testcase.value = nil
			}

			value, err := cache.Get(testcase.key)
			assert.Equal(t, testcase.expectedError, err)
			assert.Equal(t, testcase.value, value)
		})
	}

}

func TestSetValue(t *testing.T) {

	testcases := []struct {
		name  string
		key   string
		value *TestCacheValue
	}{
		{
			name: "value in the cache",
			key:  "foo",
			value: &TestCacheValue{
				value:  "foo",
				expiry: time.Now().Add(24 * time.Hour),
			},
		},
		{
			name:  "no value in the cache",
			key:   "foo",
			value: nil,
		},
		{
			name: "value in the cache but expired",
			key:  "foo",
			value: &TestCacheValue{
				value:  "foo",
				expiry: time.Now().Add(1 * time.Microsecond),
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {

			// Create new generic cache for the test value with a 24 hour cleanup duration (default)
			cache := NewGenericCache[*TestCacheValue](24 * time.Hour)

			err := cache.Set(testcase.key, testcase.value)
			assert.NoError(t, err)
			assert.Equal(t, testcase.value, cache.CacheMap[testcase.key])
		})
	}

}
