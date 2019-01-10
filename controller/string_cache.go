package controller

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type stringCache struct {
	underlying *gocache.Cache
}

func newStringCache(expiration time.Duration) *stringCache {
	c := &stringCache{
		underlying: gocache.New(expiration, expiration/2),
	}

	return c
}

func (c *stringCache) Contains(s string) bool {
	_, v := c.underlying.Get(s)
	return v
}

func (c *stringCache) Set(s string) {
	c.underlying.Set(s, struct{}{}, 0)
}
