// handlers_test.go
package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	. "sg-api/db"
	"testing"

	"github.com/matryer/is"
)

func setup(t *testing.T) (*DB, func()) {
	opts := DefaultOpts()
	db, err := NewDB(opts)
	if err != nil {
		t.Errorf("Db connection error: %s", err)
		return nil, func() {}
	}
	return db, func() {
		if err := db.Cleanup(); err != nil {
			t.Errorf("Db.Close: %s", err)
		}
	}
}

//TestHealthCheckHandler tests health check endpoint
func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	srv := NewServer()
	handler := http.HandlerFunc(srv.handleHealthCheck)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//TestHandleIndex tests index page
func TestHandleIndex(t *testing.T) {

	is := is.New(t)
	srv := NewServer()
	db, teardown := setup(t)
	defer teardown()
	srv.Db = db
	req, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	// Check the response status code is what we expect
	is.Equal(w.Code, http.StatusOK)
}

//TestNotFound tests that non existing page request returns 404
func TestNotFound(t *testing.T) {

	is := is.New(t)
	srv := NewServer()
	db, teardown := setup(t)
	defer teardown()
	srv.Db = db
	req, err := http.NewRequest("GET", "/notfound", nil)
	is.NoErr(err)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	is.Equal(w.Code, http.StatusNotFound)
}

//TestEmptyClientsList tests when on empty database empty client list is returned
func TestEmptyClientsList(t *testing.T) {

	is := is.New(t)
	srv := NewServer()
	db, teardown := setup(t)
	defer teardown()
	srv.Db = db
	req, err := http.NewRequest("GET", "/clients", nil)
	is.NoErr(err)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	is.Equal(w.Code, http.StatusOK)
	is.Equal(w.Body.String(), "[]\n")
}

//TestNewClient tests that when creating new client all the fields are properly set
func TestNewClient(t *testing.T) {

	is := is.New(t)
	srv := NewServer()
	db, teardown := setup(t)
	defer teardown()
	srv.Db = db
	var jsonStr = []byte(`{
    "salesforceId": 20401,
    "country": "Spain",
    "name": "Lexidrill",
    "owner": "Dietrich Axelsen",
    "manager": "Aristide Mullane"
  }
`)
	// create new client
	req, err := http.NewRequest("POST", "/clients", bytes.NewBuffer(jsonStr))
	is.NoErr(err)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	print(w.Body.String())
	is.Equal(w.Code, http.StatusCreated)

	// check newly created client
	req, err = http.NewRequest("GET", "/clients/20401", nil)
	is.NoErr(err)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	print(w.Body.String())
	is.Equal(w.Code, http.StatusOK)
	var client *Client
	_ = json.NewDecoder(w.Body).Decode(&client)
	is.Equal("Dietrich Axelsen", client.Owner)
	is.Equal("Aristide Mullane", client.Manager)
	is.Equal("Spain", client.Country)
	is.Equal(int32(20401), client.SalesforceID)
}
