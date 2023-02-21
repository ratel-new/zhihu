package main

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

func getCache(key string) (string, bool) {
	value, b := CacheVar.Get(key)
	if b {
		return value.(string), b
	}
	return "", b
}

func setCache(key string, value string) {
	CacheVar.Set(key, value, cache.DefaultExpiration)
}
