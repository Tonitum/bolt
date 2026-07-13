package id

import (
	"bolt/services/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterEntry(t *testing.T) {
	var alias = "foo"
	var body = IDRequest{"https://www.tonitum.com/", alias, nil}
	var bodyJson, _ = json.Marshal(body)

	r := httptest.NewRequest("GET", "/new", strings.NewReader(string(bodyJson)))
	w := httptest.NewRecorder()

	RegisterEntry(w, r)

	res := w.Result()

	if res.StatusCode != 200 {
		t.Errorf("Non 200 status code: %d", res.StatusCode)
	}

	if database.GetURL("foo") != "https://www.tonitum.com/" {
		t.Errorf("Entry not persisted to database")
	}

	var idresponse IDResponse

	resErr := json.NewDecoder(res.Body).Decode(&idresponse)

	if resErr != nil {
		t.Error("Non-JSON response")
	}

	if idresponse.ShortURL != alias {
		t.Error("Alias was not preserved")
	}

}

func TestRegisterEntryNoAlias(t *testing.T) {
	var body = IDRequest{"https://www.tonitum.com/", "", nil}
	var bodyJson, _ = json.Marshal(body)

	r := httptest.NewRequest("GET", "/new", strings.NewReader(string(bodyJson)))
	w := httptest.NewRecorder()

	RegisterEntry(w, r)

	res := w.Result()

	if res.StatusCode != 200 {
		t.Errorf("Non 200 status code: %d", res.StatusCode)
	}

	if database.GetURL("foo") != "https://www.tonitum.com/" {
		t.Errorf("Entry not persisted to database")
	}

	var idresponse IDResponse

	resErr := json.NewDecoder(res.Body).Decode(&idresponse)

	if resErr != nil {
		t.Error("Non-JSON response")
	}
}

func TestRegisterEntryNoURL(t *testing.T) {
	var body = IDRequest{"", "foo", nil}
	var bodyJson, _ = json.Marshal(body)

	r := httptest.NewRequest("GET", "/new", strings.NewReader(string(bodyJson)))
	w := httptest.NewRecorder()

	RegisterEntry(w, r)

	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Non 200 status code: %d", res.StatusCode)
	}
}
