package cache

import (
	"beedance-mcp/pkg/table"
	"time"
)

var cacheManager *CacheManager

type CacheManager struct {
	cacheTable *table.Table[CacheType, string, any]
}

func InitCacheManager() {
	mngr := &CacheManager{
		cacheTable: table.NewTable[CacheType, string, any](),
	}
	cacheManager = mngr
}

func GetByKey[V any](cacheType CacheType, key string, supplier func() any) V {
	val := cacheManager.getByKey(cacheType, key, supplier, 30)
	return val.(V)
}

func (m *CacheManager) put(cacheType CacheType, key string, val any, ttl int) {
	m.cacheTable.Put(cacheType, key, val)
	go func() {
		time.Sleep(time.Duration(ttl) * time.Second)
		m.cacheTable.Remove(cacheType, key)
	}()
}

func (m *CacheManager) getByKey(cacheType CacheType, key string, supplier func() any, ttl int) any {
	val, ok := m.cacheTable.Get(cacheType, key)
	if ok {
		return val
	}
	val = supplier()
	m.put(cacheType, key, val, ttl)
	return val
}
