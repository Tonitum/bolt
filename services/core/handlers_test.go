package core

import (
	"bolt/services/database"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockDB struct {
	data map[string]string
}

func (db *MockDB) Init() error {
	db.data = map[string]string{}
	return nil
}

func (db *MockDB) GetURL(alias string) (string, error) {
	url := db.data[alias]
	return url, nil
}

func (db *MockDB) PutAlias(alias string, url string) (bool, error) {
	db.data[alias] = url
	return true, nil
}

func (db *MockDB) DeleteAlias(alias string) (bool, error) {
	delete(db.data, alias)
	return true, nil
}

func (db *MockDB) ListAliases() (map[string]string, error) {
	return db.data, nil
}

func TestRegisterEntry(t *testing.T) {
	database.DB = &MockDB{}
	database.DB.Init()
	var alias = "foo"
	var body = AliasRequest{"https://www.tonitum.com/", alias, nil}
	var bodyJson, _ = json.Marshal(body)

	r := httptest.NewRequest("GET", "/new", strings.NewReader(string(bodyJson)))
	w := httptest.NewRecorder()

	RegisterEntry(w, r)

	res := w.Result()

	if res.StatusCode != 200 {
		t.Errorf("Non 200 status code: %d", res.StatusCode)
	}

	url, _ := database.DB.GetURL("foo")

	if url != "https://www.tonitum.com/" {
		t.Errorf("Entry not persisted to database")
	}

	var idresponse AliasResponse

	resErr := json.NewDecoder(res.Body).Decode(&idresponse)

	if resErr != nil {
		t.Error("Non-JSON response")
	}

	if idresponse.ShortURL != alias {
		t.Error("Alias was not preserved")
	}
}

func TestRegisterEntryNoAlias(t *testing.T) {
	database.DB = &MockDB{}
	database.DB.Init()
	var body = AliasRequest{"https://www.tonitum.com/", "", nil}
	var bodyJson, _ = json.Marshal(body)

	r := httptest.NewRequest("GET", "/new", strings.NewReader(string(bodyJson)))
	w := httptest.NewRecorder()

	RegisterEntry(w, r)

	res := w.Result()

	if res.StatusCode != 200 {
		t.Errorf("Non 200 status code: %d", res.StatusCode)
	}
	var idresponse AliasResponse

	resErr := json.NewDecoder(res.Body).Decode(&idresponse)

	if resErr != nil {
		t.Error("Non-JSON response")
	}

	url, _ := database.DB.GetURL(idresponse.ShortURL)

	if url != "https://www.tonitum.com/" {
		fmt.Printf("url: %v\n", url)
		t.Errorf("Entry not persisted to database")
	}
}

func TestRegisterEntryNoURL(t *testing.T) {
	var body = AliasRequest{"", "foo", nil}
	var bodyJson, _ = json.Marshal(body)

	r := httptest.NewRequest("GET", "/new", strings.NewReader(string(bodyJson)))
	w := httptest.NewRecorder()

	RegisterEntry(w, r)

	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Non 200 status code: %d", res.StatusCode)
	}
}
