package main

import (
	"fmt"
	"github.com/bernardoms/shortenerurl-go/configs"
	"github.com/bernardoms/shortenerurl-go/internal/handler"
	"github.com/bernardoms/shortenerurl-go/internal/repository"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgorilla/v1"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"os"
	"time"
)

const SHORTENERURL = "/v1/shorteners"

func main() {

	app, errNewRelic := newrelic.NewApplication(
		newrelic.NewConfig(os.Getenv("NEWRELIC_APP"), "NEWRELIC_LICENSE"),
	)

	if errNewRelic != nil {
		log.Print("Error in new relic agent")
	}

	cacheHandler := handler.CacheHandler{Cache: cache.New(5*time.Minute, 10*time.Minute) }

	shortenerHandler := handler.ShortenerHandler{Repository: initMongo(), Cache: cacheHandler}

	r := mux.NewRouter()

	r.HandleFunc(SHORTENERURL, shortenerHandler.GetAll).Methods("GET")
	r.HandleFunc(SHORTENERURL + "/{alias}", shortenerHandler.Redirect).Methods("GET")
	r.HandleFunc(SHORTENERURL, shortenerHandler.SaveShortenerURL).Methods("POST")

	nrgorilla.InstrumentRoutes(r, app)

	fmt.Printf("running server on %d", 8080)

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		fmt.Printf("error to open port %s with error %s", "8080", err)
	}
}


func initMongo() repository.Mongo{
	c := configs.NewMongoConfig()
	repository.New(c)
	mongo := repository.Mongo{Collection: repository.GetShortenerCollection()}
	return mongo
}