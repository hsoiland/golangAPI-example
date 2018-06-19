package main

import(
	"net/http"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"database/sql"
	"fmt"
	"github.com/patrickmn/go-cache"
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

func getTripsViaCache(w http.ResponseWriter, r *http.Request) {
	query:= r.URL.Query()
	medallion := query["medallion"]
	dateTime := query["dateTime"]

	log.Println(medallion)
	log.Println(dateTime)
	for i:= 0; i < len(medallion); i++ {
		medallion := medallion[i]
		dateTime := dateTime[i]
		queryStr := medallion + dateTime
		log.Println("Query is " + queryStr)
		log.Println("Checking cache first")
		result, found := c.Get(queryStr)
		log.Println(found)

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



}

func getTripsBypassCache (w http.ResponseWriter, r *http.Request) {
	query:= r.URL.Query()
	medallion := query["medallion"]
	dateTime := query["dateTime"]

	log.Println(medallion)
	log.Println(dateTime)
	for i:= 0; i < len(medallion); i++ {
		medallion := medallion[i]
		dateTime := dateTime[i]
		queryStr := medallion + dateTime
		log.Println("Query is " + queryStr)
		medallionTrips := getTripsThenCache(medallion, dateTime, queryStr)
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