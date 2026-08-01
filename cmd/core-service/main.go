package main

import (
	"bolt/internal/core"
	"bolt/internal/database"
	"log"
	"net/http"
)

func main() {
	database.InitDB()
	http.HandleFunc("POST /new", core.RegisterEntry)
	http.HandleFunc("GET /{alias}", core.LoadEntry)
	http.HandleFunc("GET /list", core.ListURLS)
	http.HandleFunc("DELETE /{alias}", core.DeleteEntry)
	log.Fatal(http.ListenAndServe(":80", nil))
}
