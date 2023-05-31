package main

import (
	"log"
	"net/http"

	"github.com/HedonisticAI/go_sitestat/http_handler"
	"github.com/gorilla/mux"

	"github.com/HedonisticAI/go_sitestat/db_handler"
)

func main() {
	db, err := db_handler.InitPostgres()
	if err != nil {
		log.Fatal(err)
		return
	}
	redClient := db_handler.InitRedis()
	collection := &http_handler.DBInstance{redClient, db}
	go func() {
		http_handler.GetTime(collection)
	}()
	r := mux.NewRouter()
	r.HandleFunc("/site/{sitename}", collection.Site)
	r.HandleFunc("/max", collection.Max)
	r.HandleFunc("/min", collection.Min)
	r.HandleFunc("/admsite/{sitename}", collection.AdmSite)
	r.HandleFunc("/admMin", collection.AdmMin)
	r.HandleFunc("/admMax", collection.AdmMin)
	http.ListenAndServe(":8000", r)
	defer redClient.Close()
}
