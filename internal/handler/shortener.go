package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/bernardoms/shortenerurl-go/internal/cache"
	"github.com/bernardoms/shortenerurl-go/internal/repository"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type ShortenerHandler struct {
	Repository repository.Repository
	Cache cache.Cache
}

func (s ShortenerHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	result, err := s.Repository.FindAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJson(w, http.StatusOK, result, "")
}

func (s ShortenerHandler) Redirect(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	cached, found := s.Cache.GetCache(vars["alias"])

	var result *repository.URLShortener

	if found {
		result = cached.(*repository.URLShortener)
		result.RedirectCount = result.RedirectCount + 1
		updated, err := s.Repository.Update(result)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusFound, "", updated.OriginalUrl)
	} else {
		result, err := s.Repository.FindByAlias(vars["alias"])
		if err != nil {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			s.Cache.SetCache(vars["alias"], result, 30 * time.Minute)
			result.RedirectCount = result.RedirectCount + 1
			updated, err := s.Repository.Update(result)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondWithJson(w, http.StatusFound, "", updated.OriginalUrl)
		}
	}
}

func (s ShortenerHandler) SaveShortenerURL(w http.ResponseWriter, r *http.Request)  {
	var shortener *repository.URLShortener

	err := json.NewDecoder(r.Body).Decode(&shortener)

	shortener.Id = primitive.NewObjectID()
	shortener.Alias = generateRandomAlias()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	inserted, err := s.Repository.Save(shortener)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {

		respondWithJson(w, http.StatusCreated, "", "v1/shorteners/" + inserted.Alias)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg}, "")
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}, location string) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")

	if location != "" {
		w.Header().Set("Location", location)
		w.WriteHeader(code)
		return
	}
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func generateRandomAlias () string {
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}