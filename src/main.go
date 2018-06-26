package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"github.com/patrickmn/go-cache"
)

//"D7D598CD99978BD012A87A76A7C891B7"
//"2013-12-01 00:13:00"

var (
	c = cache.New(cache.NoExpiration, cache.NoExpiration)
)

func init() {
	log.Println("testing cache")
	c.Set("foo", "bar", cache.NoExpiration)
	foo, found := c.Get("foo")
	if found {
		log.Println(foo)
	}
	bar, found := c.Get("bar")
	if found {
		log.Println("cache doesnt work properly")
	} else {
		log.Println(bar)
		log.Println("bar not found, cache working")
	}
}

func main() {
	var router = mux.NewRouter()

	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/getTrips", multipleDatesViaCache).Methods("GET")
	router.HandleFunc("/getTrips/bypassCache", multipleDatesBypassCache).Methods("GET")
	router.HandleFunc("/getTrips/singleDate", singleDateViaCache).Methods("GET")
	router.HandleFunc("/getTrips/singleDate/bypassCache", singleDateBypassCache).Methods("GET")
	router.HandleFunc("/flushCache", clearCache).Methods("GET")
	log.Printf("CabTrips is running on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}

func clearCache (w http.ResponseWriter, r *http.Request) {
	log.Println("Flushing Cache")
	c.Set("foo", "bar", cache.NoExpiration)
	json.NewEncoder(w).Encode("Flushing Cache")
	c.Flush()
	foo, found :=c.Get("foo")
	if found {
		log.Println(foo)
		log.Fatal("not flushed")
	}
}
