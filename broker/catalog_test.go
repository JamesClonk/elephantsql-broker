package broker

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JamesClonk/elephantsql-broker/log"
	"github.com/JamesClonk/elephantsql-broker/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestBroker_Catalog(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/regions": util.HttpTestCase{200, util.Body("../_fixtures/api_list_regions.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/catalog", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id": "8ff5d1c8-c6eb-4f04-928c-6a422e0ea330"`)
	assert.Contains(t, rec.Body.String(), `"displayName": "ElephantSQL"`)
	assert.Contains(t, rec.Body.String(), `"name": "turtle"`)
	assert.Contains(t, rec.Body.String(), `"imageUrl": "https://www.elephantsql.com/images/turtle_256.png"`)
	assert.Contains(t, rec.Body.String(), `"imageUrl": "https://www.elephantsql.com/images/dolphin_256.png"`)
	assert.Contains(t, rec.Body.String(), `"Shared high performance server"`)
	assert.Contains(t, rec.Body.String(), `"Dedicated server"`)
	assert.Contains(t, rec.Body.String(), `"description": "Ruthless Rat - dedicated instance with follower"`)
	assert.Equal(t, util.Body("../_fixtures/broker_catalog.json"), rec.Body.String())
}

func TestBroker_Catalog_NoSharedPlans(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/regions": util.HttpTestCase{200, util.Body("../_fixtures/api_list_regions.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()

	testConfig := util.TestConfig(apiServer.URL)
	testConfig.API.DefaultRegion = "google-compute-engine::europe-west6" // region does not have shared plans
	r := NewRouter(testConfig)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/catalog", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id": "8ff5d1c8-c6eb-4f04-928c-6a422e0ea330"`)
	assert.Contains(t, rec.Body.String(), `"displayName": "ElephantSQL"`)
	assert.NotContains(t, rec.Body.String(), `"name": "turtle"`)
	assert.NotContains(t, rec.Body.String(), `"imageUrl": "https://www.elephantsql.com/images/turtle_256.png"`)
	assert.Contains(t, rec.Body.String(), `"imageUrl": "https://www.elephantsql.com/images/dolphin_256.png"`)
	assert.NotContains(t, rec.Body.String(), `"Shared high performance server"`)
	assert.Contains(t, rec.Body.String(), `"Dedicated server"`)
	assert.Contains(t, rec.Body.String(), `"description": "Ruthless Rat - dedicated instance with follower"`)
	assert.NotEqual(t, util.Body("../_fixtures/broker_catalog.json"), rec.Body.String())
}
