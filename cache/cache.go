package cache

import (
	"sync"
	"time"
)

var (
	cacheData   map[string]interface{}
	cacheMutex  sync.RWMutex
	defaultCacheDuration = 20 * time.Second
	lastUpdated time.Time
)

func IsExpired() bool {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return time.Since(lastUpdated) > defaultCacheDuration
}

func GetLastUpdated() time.Time {
	return lastUpdated
}

func GetCache() map[string]interface{} {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	// Return a copy of the cache
	// cacheCopy := make(map[string]interface{})
	// for k, v := range cacheData {
	// 	cacheCopy[k] = v
	// }

	return cacheData
}

func SetCache(data map[string]interface{}) {
	cacheMutex.Lock()  // Lock for writing
	defer cacheMutex.Unlock()

	cacheData = data
	lastUpdated = time.Now()
	
}
