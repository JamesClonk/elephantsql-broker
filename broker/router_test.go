package broker

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/JamesClonk/elephantsql-broker/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Chdir("../")
}

func TestBroker_HealthEndpoint(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	NewRouter(util.TestConfig("")).ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, `{ "status": "ok" }`, rec.Body.String())
}

func TestBroker_BasicAuthValid(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")

	NewRouter(util.TestConfig("")).ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, `{ "status": "ok" }`, rec.Body.String())
}

func TestBroker_BasicAuthUnauthorized(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	NewRouter(util.TestConfig("")).ServeHTTP(rec, req)
	assert.Equal(t, 401, rec.Code)
	assert.Equal(t, `{ "error": "Unauthorized", "description": "You are not authorized to access this service broker" }`, rec.Body.String())
}
