package main

import(
	"net/http"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"database/sql"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/gorilla/mux"
)

type Trip struct {
	ID string `json:"id"`
	Time string `json:"time"`
	Count int `json:"count"`
}

type CountTrips struct {
	ID string `json:"id"`
	Count int `json:"count"`
}

func singleDateViaCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	medallion := vars["medallion"]
	dateTime := vars["dateTime"]
	getTripsViaCache(medallion, dateTime,w)
}

func multipleDatesViaCache(w http.ResponseWriter, r *http.Request) {
	query:= r.URL.Query()
	medallion := query["medallion"]
	dateTime := query["dateTime"]
	for i:= 0; i < len(medallion); i++ {
		medallion := medallion[i]
		dateTime := dateTime[i]
		getTripsViaCache(medallion, dateTime, w)
	}
}

func singleDateBypassCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	medallion := vars["medallion"]
	dateTime := vars["dateTime"]
	queryStr := medallion + dateTime
	log.Println("Query is " + queryStr)
	medallionTrips := getTripsThenCache(medallion, dateTime, queryStr)
	json.NewEncoder(w).Encode(medallionTrips)
}

func multipleDatesBypassCache(w http.ResponseWriter, r *http.Request) {
	query:= r.URL.Query()
	medallion := query["medallion"]
	dateTime := query["dateTime"]
	for i:= 0; i < len(medallion); i++ {
		medallion := medallion[i]
		dateTime := dateTime[i]
		queryStr := medallion + dateTime
		log.Println("Query is " + queryStr)
		medallionTrips := getTripsThenCache(medallion, dateTime, queryStr)
		json.NewEncoder(w).Encode(medallionTrips)
	}
}

func getTripsViaCache(medallion string, dateTime string, w http.ResponseWriter) {
	queryStr := medallion + dateTime
	log.Println("Query is " + queryStr)
	log.Println("Checking cache first")
	result, found := c.Get(queryStr)
		if found {
			log.Printf("Cache returned result")
			var medallionTrips CountTrips
			medallionTrips.ID = medallion
			medallionTrips.Count = result.(int)
			json.NewEncoder(w).Encode(medallionTrips)
			log.Println(medallionTrips)

		} else {
			log.Println("Does not exist in cache")
			medallionTrips:= getTripsThenCache(medallion, dateTime, queryStr)
			json.NewEncoder(w).Encode(medallionTrips)

		}

}


func getTripsThenCache(medallion string, dateTime string, query string) (CountTrips) {
	count := dbQuery(medallion, dateTime)
	c.Set(query, count, cache.NoExpiration)
	var medallionTrips CountTrips
	medallionTrips.ID = medallion
	medallionTrips.Count = count
	log.Println(medallionTrips)
	return medallionTrips

}

func dbQuery(medallion string, dateTime string) (int){
	log.Println("Opening SQL connection")
	// note for production the password should be set in the server for access by the application in a non-vulnerable fashion
	// and passwords should never be stored in code
	db, err := sql.Open("mysql", "root:root@/cabdata")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Value not in cache querying database")
	var sqlquery = fmt.Sprintf("SELECT medallion, pickup_datetime FROM `cabdata`.`cab_trip_data` WHERE medallion = '%s' AND pickup_datetime = '%s'", medallion, dateTime)
	results, err := db.Query(sqlquery)
	if err != nil {
		log.Fatal("Query not valid")
	}
	count := 0
	for results.Next() {
		var trip Trip
		err = results.Scan(&trip.ID, &trip.Time)
		if err != nil {
			log.Fatal(err.Error())
		}
		count++
	}
	defer db.Close()

	return count
}