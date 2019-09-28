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

func TestBroker_FetchServiceInstance(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_fetch_service_instance.json"), rec.Body.String())
}

func TestBroker_FetchServiceInstanceNotFound(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances": util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/52551d5f-1350-4f7d-9ddd-710a47ef9b72", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ServiceInstanceNotFound"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}
