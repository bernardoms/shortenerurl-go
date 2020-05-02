package mock

import (
	"github.com/bernardoms/shortenerurl-go/internal/repository"
	"github.com/stretchr/testify/mock"
)

type MongoMock struct {
	mock.Mock
}

func (m *MongoMock) FindAll() ([] *repository.URLShortener, error) {
	args := m.Called()
	return args.Get(0).([] *repository.URLShortener), args.Error(1)
}

func (m *MongoMock) FindByAlias(alias string) (*repository.URLShortener, error) {
	args := m.Called(alias)
	return args.Get(0).(*repository.URLShortener), args.Error(1)
}

func (m *MongoMock) Save(shortener *repository.URLShortener) (*repository.URLShortener, error) {
	args := m.Called(shortener)
	return args.Get(0).(*repository.URLShortener), args.Error(1)
}

func (m *MongoMock) Update (shortener *repository.URLShortener) (*repository.URLShortener, error) {
	args := m.Called(shortener)
	return args.Get(0).(*repository.URLShortener), args.Error(1)
}