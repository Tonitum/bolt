package main

import (
	"bolt/services/database"
	"bolt/services/id"
	"log"
	"net/http"
)

func main() {
	database.InitDB()
	http.HandleFunc("POST /new", id.RegisterEntry)
	http.HandleFunc("GET /{long_url}",id.LoadURL)
	http.HandleFunc("GET /list",id.ListURLS)
	log.Fatal(http.ListenAndServe(":80", nil))
}
