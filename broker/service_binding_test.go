package broker

import (
	"bytes"
	"encoding/json"
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

func TestBroker_BindServiceBinding(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	binding := ServiceBinding{
		ServiceID: "8ff5d1c8-c6eb-4f04-928c-6a422e0ea330",
		PlanID:    "890d1ed6-0ff6-4c93-afc6-df753be6f1e3",
	}
	data, _ := json.MarshalIndent(binding, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_fetch_service_binding.json"), rec.Body.String())
}

func TestBroker_BindServiceBinding_InstanceNotFound(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_BindServiceBinding_EmptyBody(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not read binding request"`)
}

func TestBroker_BindServiceBinding_UnknownService(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	binding := ServiceBinding{
		ServiceID: "e5356e96-2a4c-46a0-b554-81c09eef9234",
		PlanID:    "890d1ed6-0ff6-4c93-afc6-df753be6f1e3",
	}
	data, _ := json.MarshalIndent(binding, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Unknown service_id"`)
}

func TestBroker_BindServiceBinding_UnknownPlan(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	binding := ServiceBinding{
		ServiceID: "8ff5d1c8-c6eb-4f04-928c-6a422e0ea330",
		PlanID:    "e5356e96-2a4c-46a0-b554-81c09eef9234",
	}
	data, _ := json.MarshalIndent(binding, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Unknown plan_id"`)
}

func TestBroker_FetchServiceBinding(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_fetch_service_binding.json"), rec.Body.String())
}

func TestBroker_FetchServiceBinding_NotFound(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances": util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v2/service_instances/52551d5f-1350-4f7d-9ddd-710a47ef9b72/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "ServiceInstanceNotFound"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_UnbindServiceBinding(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "{}", rec.Body.String())
}

func TestBroker_UnbindServiceBinding_NotFound(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances": util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/52551d5f-1350-4f7d-9ddd-710a47ef9b72/service_bindings/deadbeef", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}
