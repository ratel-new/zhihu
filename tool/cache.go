package tool

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	CacheVar *cache.Cache
)

func init() {
	CacheVar = cache.New(5*time.Minute, 10*time.Minute)
}

func GetHtmlCache(key string) (*Fill, bool) {
	value, b := CacheVar.Get(key)
	if b {
		return value.(*Fill), b
	}
	return nil, b
}

func SetCache(key string, value *Fill) {
	CacheVar.Set(key, value, cache.DefaultExpiration)
}
