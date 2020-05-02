package handler
import (
	"github.com/patrickmn/go-cache"
	"time"
)


type CacheHandler struct {
	Cache *cache.Cache
}

 func (c CacheHandler) GetCache (key string) (interface{}, bool) {
	return c.Cache.Get(key)
}

func (c CacheHandler ) SetCache (key string, value interface{}, time time.Duration) {
	c.Cache.Set(key, value, time)
}