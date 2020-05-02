package handler

import (
	"bytes"
	"errors"
	"github.com/bernardoms/shortenerurl-go/internal/handler"
	"github.com/bernardoms/shortenerurl-go/internal/repository"
	"github.com/bernardoms/shortenerurl-go/test/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAllHandlerSuccess(t *testing.T) {

	mongoMock := mock.MongoMock{}

	h := handler.ShortenerHandler{Repository: &mongoMock}

	r, _ := http.NewRequest("GET", "/v1/shorteners", nil)
	w := httptest.NewRecorder()

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")
	id2, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994eda")

	var shorteners = [] * repository.URLShortener {
		{Id: id, OriginalUrl: "http://www.test.com", RedirectCount: 0, Alias: "test"},
		{Id: id2, OriginalUrl: "http://www.test2.com", RedirectCount: 0, Alias: "test2"},
	}

	mongoMock.On("FindAll").Return(shorteners, nil)
	h.GetAll(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"id\":\"5ea7208049e00ddb76994ede\",\"originalURL\":\"http://www.test.com\",\"alias\":\"test\",\"redirectCount\":0},{\"id\":\"5ea7208049e00ddb76994eda\",\"originalURL\":\"http://www.test2.com\",\"alias\":\"test2\",\"redirectCount\":0}]", w.Body.String())
}

func TestGetAllHandlerMongoFail(t *testing.T) {

	mongoMock := mock.MongoMock{}

	shortenerHandler := handler.ShortenerHandler{Repository: &mongoMock}

	r, _ := http.NewRequest("GET", "/v1/shorteners", nil)
	w := httptest.NewRecorder()

	mongoMock.On("FindAll").Return([]*repository.URLShortener{}, errors.New("error on mongo"))

	shortenerHandler.GetAll(w, r)

	mongoMock.AssertExpectations(t)

	mongoMock.AssertNumberOfCalls(t, "FindAll", 1)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRedirectWithSuccessNoCache(t *testing.T) {
	mongoMock := mock.MongoMock{}

	cacheMock := mock.CacheHandlerMock{}

	shortenerHandler := handler.ShortenerHandler{Repository: &mongoMock, Cache: &cacheMock}

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	r, _ := http.NewRequest("GET", "/v1/shorteners/123", nil)

	w := httptest.NewRecorder()

	vars := map[string]string{
		"alias": "123",
	}

	r = mux.SetURLVars(r, vars)

	shortener  := &repository.URLShortener {Id: id, OriginalUrl: "http://www.test2.com", RedirectCount: 0, Alias: "123"}

	mongoMock.On("FindByAlias", "123").Return(shortener, nil)

	mongoMock.On("Update", shortener).Return(shortener, nil)

	cacheMock.On("GetCache", "123").Return("http://www.test2.com", false)

	cacheMock.On("SetCache", "123", shortener, 30 * time.Minute).Return(nil)

	shortenerHandler.Redirect(w, r)

	mongoMock.AssertExpectations(t)

	mongoMock.AssertNumberOfCalls(t, "FindByAlias", 1)

	mongoMock.AssertNumberOfCalls(t, "Update", 1)

	cacheMock.AssertNumberOfCalls(t, "GetCache", 1)

	cacheMock.AssertNumberOfCalls(t, "SetCache", 1)

	assert.Equal(t, http.StatusFound, w.Code)
}

func TestRedirectWithSuccessCached (t *testing.T) {
	mongoMock := mock.MongoMock{}

	cacheMock := mock.CacheHandlerMock{}

	shortenerHandler := handler.ShortenerHandler{Repository: &mongoMock, Cache: &cacheMock}

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	r, _ := http.NewRequest("GET", "/v1/shorteners/123", nil)

	w := httptest.NewRecorder()

	vars := map[string]string{
		"alias": "123",
	}

	r = mux.SetURLVars(r, vars)

	shortener  := &repository.URLShortener {Id: id, OriginalUrl: "http://www.test2.com", RedirectCount: 0, Alias: "123"}

	mongoMock.On("Update", shortener).Return(shortener, nil)

	cacheMock.On("GetCache", "123").Return(shortener, true)

	shortenerHandler.Redirect(w, r)

	mongoMock.AssertExpectations(t)

	mongoMock.AssertNotCalled(t, "FindByAlias")

	mongoMock.AssertNumberOfCalls(t, "Update", 1)

	cacheMock.AssertNumberOfCalls(t, "GetCache", 1)

	cacheMock.AssertNotCalled(t, "SetCache")

	assert.Equal(t, http.StatusFound, w.Code)
}

func TestRedirectWithFailOnFind(t *testing.T) {
	mongoMock := mock.MongoMock{}

	cacheMock := mock.CacheHandlerMock{}

	shortenerHandler := handler.ShortenerHandler{Repository: &mongoMock, Cache: &cacheMock}

	r, _ := http.NewRequest("GET", "/v1/shorteners/123", nil)

	w := httptest.NewRecorder()

	vars := map[string]string{
		"alias": "123",
	}

	r = mux.SetURLVars(r, vars)

	mongoMock.On("FindByAlias", "123").Return(&repository.URLShortener{}, errors.New("document not found"))

	cacheMock.On("GetCache", "123").Return(nil, false)

	shortenerHandler.Redirect(w, r)

	mongoMock.AssertExpectations(t)

	mongoMock.AssertNumberOfCalls(t, "FindByAlias", 1)

	mongoMock.AssertNotCalled(t, "Update")

	assert.Equal(t, http.StatusNotFound, w.Code)
}


func TestRedirectWithFailOnUpdate(t *testing.T) {
	mongoMock := mock.MongoMock{}

	cacheMock := mock.CacheHandlerMock{}

	shortenerHandler := handler.ShortenerHandler{Repository: &mongoMock, Cache: &cacheMock}

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	r, _ := http.NewRequest("GET", "/v1/shorteners/123", nil)

	w := httptest.NewRecorder()

	vars := map[string]string{
		"alias": "123",
	}

	r = mux.SetURLVars(r, vars)

	shortener  := &repository.URLShortener {Id: id, OriginalUrl: "http://www.test2.com", RedirectCount: 0, Alias: "test2"}

	mongoMock.On("FindByAlias", "123").Return(shortener, nil)

	mongoMock.On("Update", shortener).Return(&repository.URLShortener{}, errors.New("fail to update"))

	cacheMock.On("GetCache", "123").Return(nil, false)

	cacheMock.On("SetCache", "123", shortener, 30 * time.Minute).Return(nil)

	shortenerHandler.Redirect(w, r)

	mongoMock.AssertExpectations(t)

	mongoMock.AssertNumberOfCalls(t, "FindByAlias", 1)

	mongoMock.AssertNumberOfCalls(t, "Update", 1)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSaveShortener(t *testing.T) {

	mongoMock := mock.MongoMock{}

	shortenerHandler := handler.ShortenerHandler{Repository: &mongoMock}

	jsonStr := []byte(`{"originalUrl":"http://www.test2.com"}`)

	r, _ := http.NewRequest("POST", "/v1/shorteners", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	shortener  := &repository.URLShortener {Id: id, OriginalUrl: "http://www.test2.com", RedirectCount: 0, Alias: "test2"}

	mongoMock.On("Save", mock2.Anything).Return(shortener, nil)

	shortenerHandler.SaveShortenerURL(w, r)

	mongoMock.AssertExpectations(t)

	mongoMock.AssertNumberOfCalls(t, "Save", 1)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "v1/shorteners/test2", w.Header().Get("Location"))
}

func TestSaveShortenerError(t *testing.T) {

	mongoMock := mock.MongoMock{}

	shortenerHandler := handler.ShortenerHandler{Repository: &mongoMock}

	jsonStr := []byte(`{"originalUrl":"http://www.test2.com"}`)

	r, _ := http.NewRequest("POST", "/v1/shorteners", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()

	mongoMock.On("Save", mock2.Anything).Return(&repository.URLShortener{}, errors.New("ERROR"))

	shortenerHandler.SaveShortenerURL(w, r)

	mongoMock.AssertExpectations(t)

	mongoMock.AssertNumberOfCalls(t, "Save", 1)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "", w.Header().Get("Location"))
}