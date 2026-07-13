package main

import (
	"bolt/services/id"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("POST /new", id.RegisterEntry)
	http.HandleFunc("GET /{long_url}",id.LoadURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
