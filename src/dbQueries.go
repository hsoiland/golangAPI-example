package main

import(
	"net/http"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/patrickmn/go-cache"
	"encoding/json"
	"log"
)



type Trip struct {
	ID string `json:"id"`
	Time string `json:"time"`
	Count int `json:"count"`
}

type Trips []Trip

func singleDateViaCache(w http.ResponseWriter, r *http.Request) {
	query:= r.URL.Query()
	medallion := query["medallion"]
	dateTime := query["dateTime"]
	var trips Trips
	for i:= 0; i < len(medallion); i++ {
		medallion := medallion[i]
		medallionTrip := getTripsViaCache(medallion, dateTime[0])
		trips = append(trips, medallionTrip)
	}
	json.NewEncoder(w).Encode(trips)
	log.Println(trips)
}

func multipleDatesViaCache(w http.ResponseWriter, r *http.Request) {
	query:= r.URL.Query()
	medallion := query["medallion"]
	dateTime := query["dateTime"]
	var trips Trips
	for i:= 0; i < len(medallion); i++ {
		medallion := medallion[i]
		dateTime := dateTime[i]
		medallionTrip := getTripsViaCache(medallion, dateTime)
		trips = append(trips, medallionTrip)
	}
	json.NewEncoder(w).Encode(trips)
	log.Println(trips)
}

func singleDateBypassCache(w http.ResponseWriter, r *http.Request) {
	query:= r.URL.Query()
	medallion := query["medallion"]
	dateTime := query["dateTime"]
	var trips Trips
	for i:= 0; i < len(medallion); i++ {
		medallion := medallion[i]
		queryStr := medallion + dateTime[0]
		fmt.Println("Query is " + queryStr)
		medallionTrip := getTripThenCache(medallion, dateTime[0], queryStr)
		trips = append(trips, medallionTrip)
	}
	json.NewEncoder(w).Encode(trips)
	log.Println(trips)
}

func multipleDatesBypassCache(w http.ResponseWriter, r *http.Request) {
	query:= r.URL.Query()
	medallion := query["medallion"]
	dateTime := query["dateTime"]
	var trips Trips
	for i:= 0; i < len(medallion); i++ {
		medallion := medallion[i]
		dateTime := dateTime[i]
		queryStr := medallion + dateTime
		fmt.Println("Query is " + queryStr)
		medallionTrip := getTripThenCache(medallion, dateTime, queryStr)
		trips = append(trips, medallionTrip)
	}
	json.NewEncoder(w).Encode(trips)
	log.Println(trips)
}

func getTripsViaCache(medallion string, dateTime string) (Trip) {
	queryStr := medallion + dateTime
	fmt.Println("Query is " + queryStr)
	fmt.Println("Checking cache first")
	result, found := c.Get(queryStr)
		if found {
			fmt.Println("Cache returned result")
			var medallionTrip Trip
			medallionTrip.ID = medallion
			medallionTrip.Time = dateTime
			medallionTrip.Count = result.(int)
			return medallionTrip
		} else {
			fmt.Println("Does not exist in cache")
			medallionTrip:= getTripThenCache(medallion, dateTime, queryStr)
			return medallionTrip
		}

}


func getTripThenCache(medallion string, dateTime string, query string) (Trip) {
	count := dbQuery(medallion, dateTime)
	c.Set(query, count, cache.NoExpiration)
	var medallionTrip Trip
	medallionTrip.ID = medallion
	medallionTrip.Time = dateTime
	medallionTrip.Count = count
	return medallionTrip
}

func dbQuery(medallion string, dateTime string) (int){
	fmt.Println("Opening SQL connection")
	// note for production the password should be set in the server for access by the application in a non-vulnerable fashion
	// and passwords should never be stored in code
	db, err := sql.Open("mysql", "root:root@/cabdata")
	if err != nil {
		log.Fatal(err.Error())
	}
	var sqlquery = fmt.Sprintf("SELECT medallion, pickup_datetime FROM `cabdata`.`cab_trip_data` WHERE medallion = '%s' AND pickup_datetime = '%s'", medallion, dateTime)
	results, err := db.Query(sqlquery)
	if err != nil {
		log.Fatal("Query not valid")
	}
	count := 0
	for results.Next() {
		if err != nil {
			log.Fatal(err.Error())
		}
		count++
	}
	defer db.Close()
	return count
}