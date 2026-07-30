package core

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

type AliasRequest struct {
	Url   string `json:"url"`
	Alias string `json:"alias"` // make optional
	User  *User
}

type AliasResponse struct {
	ShortURL string `json:"short_url"`
}

func RegisterEntry(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	// authenticate the request TODO

	// parse the provided long url
	req.Body = http.MaxBytesReader(w, req.Body, 1048576)

	var aliasrequest AliasRequest
	err := json.NewDecoder(req.Body).Decode(&aliasrequest)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	fmt.Printf("received request for url %v\n", aliasrequest.Url)
	// check if the body included a short URL alias
	var reqAlias = aliasrequest.Alias
	var shortAlias string = reqAlias
	// create the shortened url
	if reqAlias == "" {
		shortAlias = getShortAlias()
	}
	// save the short URL to the database as shortURL:longURL
	if aliasrequest.Url == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	database.DB.PutAlias(shortAlias, aliasrequest.Url)

	// return the payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := AliasResponse{shortAlias}
	json.NewEncoder(w).Encode(response)
}

func getShortAlias() string {
	// - generate short alias for url
	//   - alias: group identifier (region), local identifier, random hash + incrementing counter
	//   - ex. cluster1-0-ijanek-0 or similar
	//   - base64 hash? nice for obfuscation, but not necessary maybe?
	shortAlias := "foo"  // TODO: Read from env var. Should be unique per deployment
	shortAlias += "-bar" // TODO: Read from env var. Should be unique per instance
	// append and increment global counter
	shortAlias += "-" + strconv.FormatInt(count.Add(1)-1, 10)
	// Generate 8 random characters
	shortAlias += rand.Text()[0:8] // only first 8 characters
	encoded := base64.RawURLEncoding.EncodeToString([]byte(shortAlias))

	return encoded
}

func LoadEntry(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	// authenticate the request TODO
	fullURL, err := database.DB.GetURL(strings.Split(req.URL.Path, "/")[1])
	if err != nil {
		// TODO: this is only one of the possible error cases
		http.Error(w, "alias does not exist", http.StatusNotFound)
		return
	}
	println("redirecting to " + fullURL)

	// get the short alias from the loaded url
	// get the url from the stored alias map
	// redirect the request to the long url
	http.Redirect(w, req, fullURL, http.StatusFound)
}

func ListURLS(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// return the payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	results, err := database.DB.ListAliases()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}
