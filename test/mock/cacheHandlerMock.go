package mock

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type CacheHandlerMock struct {
	mock.Mock
}

func (m *CacheHandlerMock) GetCache (key string) (interface{}, bool)  {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *CacheHandlerMock) SetCache (key string, value interface{}, time time.Duration) {
	m.Called(key, value, time)
}