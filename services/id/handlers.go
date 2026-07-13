package id

import (
	"bolt/services/database"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

var count atomic.Int64

type User struct {
	Id string `json:"id"`
}

type IDRequest struct {
	Url  string `json:"url"`
	Alias   string `json:"alias"` // make optional
	User *User
}

type IDResponse struct {
	ShortURL string `json:"short_url"`
}

func RegisterEntry(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	// authenticate the request TODO

	// parse the provided long url
	req.Body = http.MaxBytesReader(w, req.Body, 1048576)

	var idrequest IDRequest
	err := json.NewDecoder(req.Body).Decode(&idrequest)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	fmt.Printf("received request for url %v\n", idrequest.Url)
	// check if the body included a short URL id
	var shortID string
	var reqID = idrequest.Alias
	// create the shortened url
	if reqID == "" {
		shortID = getShortID()
	} else {
		shortID = reqID
	}

	// save the short URL to the database as shortURL:longURL
	if (idrequest.Url == "") {
		w.WriteHeader(http.StatusBadRequest)
	}
	database.PutID(shortID, idrequest.Url)

	// return the payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := IDResponse{shortID}
	json.NewEncoder(w).Encode(response)
}

func getShortID() string {
	// - generate short id for url
	//   - id: group identifier (region), local identifier, random hash + incrementing counter
	//   - ex. cluster1-0-ijanek-0 or similar
	//   - base64 hash? nice for obfuscation, but not necessary maybe?
	shortID := "foo"  // TODO: Read from env var. Should be unique per deployment
	shortID += "-bar" // TODO: Read from env var. Should be unique per instance
	// append and increment global counter
	shortID += "-" + strconv.FormatInt(count.Add(1)-1, 10)
	// Generate 8 random characters
	shortID += rand.Text()[0:8] // only first 8 characters
	encoded := base64.RawURLEncoding.EncodeToString([]byte(shortID))

	return encoded
}

func LoadURL(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	// authenticate the request TODO
	var fullURL string = database.GetURL(strings.Split(req.URL.Path, "/")[1])
	println("redirecting to " + fullURL)

	// get the short id from the loaded url
	// get the url from the stored id map
	// redirect the request to the long url
	http.Redirect(w, req, fullURL, http.StatusFound)
}

func ListURLS(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// return the payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(database.Dump())
}
