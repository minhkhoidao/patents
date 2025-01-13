package cache

import "time"

type CacheConfig struct {
	TTL       time.Duration
	KeyPrefix string
}

func NewCacheConfig() *CacheConfig {
	return &CacheConfig{
		TTL:       time.Hour,
		KeyPrefix: "patents:",
	}
}
