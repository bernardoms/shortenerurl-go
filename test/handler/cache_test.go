package handler

import (
	"github.com/bernardoms/shortenerurl-go/internal/handler"
	"github.com/bernardoms/shortenerurl-go/internal/repository"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSetAndGetCache(t *testing.T) {
	cacheHandler := handler.CacheHandler{Cache: cache.New(5*time.Minute, 10*time.Minute)}
	cacheHandler.SetCache("123", repository.URLShortener{}, time.Minute)
	cached, found := cacheHandler.GetCache("123")
	assert.Equal(t, found, true)
	assert.Equal(t, cached,  repository.URLShortener{})
}

func TestGetCacheWithoutSetting(t *testing.T) {
	cacheHandler := handler.CacheHandler{Cache: cache.New(5*time.Minute, 10*time.Minute)}
	cached, found := cacheHandler.GetCache("123")
	assert.Equal(t, found, false)
	assert.Equal(t, cached, nil)
}
