package main

import (
	"bolt/services/database"
	"bolt/services/core"
	"log"
	"net/http"
)

func main() {
	database.InitDB()
	http.HandleFunc("POST /new", core.RegisterEntry)
	http.HandleFunc("GET /{alias}",core.LoadEntry)
	http.HandleFunc("GET /list",core.ListURLS)
	log.Fatal(http.ListenAndServe(":80", nil))
}
