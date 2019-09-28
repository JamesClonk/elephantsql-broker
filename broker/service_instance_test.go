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

func TestBroker_ProvisionServiceInstance(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"GET::/instances":  util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"POST::/instances": util.HttpTestCase{200, util.Body("../_fixtures/api_create_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "8ff5d1c8-c6eb-4f04-928c-6a422e0ea330",
		PlanID:    "890d1ed6-0ff6-4c93-afc6-df753be6f1e3",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/47d80867-a199-4b5c-8425-8caec83151a4", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 201, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_provision_service_instance.json"), rec.Body.String())
}

func TestBroker_ProvisionServiceInstance_EmptyBody(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not read provisioning request"`)
}

func TestBroker_ProvisionServiceInstance_UnknownPlan(t *testing.T) {
	r := NewRouter(util.TestConfig(""))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "8ff5d1c8-c6eb-4f04-928c-6a422e0ea330",
		PlanID:    "deadbeef",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 400, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MalformedRequest"`)
	assert.Contains(t, rec.Body.String(), `"description": "Unknown plan_id"`)
}

func TestBroker_ProvisionServiceInstance_Conflicting(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "8ff5d1c8-c6eb-4f04-928c-6a422e0ea330",
		PlanID:    "890d1ed6-0ff6-4c93-afc6-df753be6f1e3",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 409, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "Conflict"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance already exists with different attributes"`)
}

func TestBroker_ProvisionServiceInstance_Existing(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":      util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"/instances/4567": util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	provisioning := ServiceInstanceProvisioning{
		ServiceID: "8ff5d1c8-c6eb-4f04-928c-6a422e0ea330",
		PlanID:    "6203b8e7-9ef4-44ef-bb0b-48b50409794d",
	}
	data, _ := json.MarshalIndent(provisioning, "", "  ")

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, util.Body("../_fixtures/broker_provision_service_instance.json"), rec.Body.String())
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

func TestBroker_FetchServiceInstance_NotFound(t *testing.T) {
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

func TestBroker_DeprovisionServiceInstance(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":              util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"GET::/instances/4567":    util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
		"DELETE::/instances/4567": util.HttpTestCase{204, "", nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "{}", rec.Body.String())
}

func TestBroker_DeprovisionServiceInstance_NotFound(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances": util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/52551d5f-1350-4f7d-9ddd-710a47ef9b72", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 410, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "MissingServiceInstance"`)
	assert.Contains(t, rec.Body.String(), `"description": "The service instance does not exist"`)
}

func TestBroker_DeprovisionServiceInstance_Error(t *testing.T) {
	test := map[string]util.HttpTestCase{
		"/instances":              util.HttpTestCase{200, util.Body("../_fixtures/api_list_instances.json"), nil},
		"GET::/instances/4567":    util.HttpTestCase{200, util.Body("../_fixtures/api_get_instance.json"), nil},
		"DELETE::/instances/4567": util.HttpTestCase{404, "", nil},
	}
	apiServer := util.TestServer("", "deadbeef", test)
	defer apiServer.Close()
	r := NewRouter(util.TestConfig(apiServer.URL))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v2/service_instances/8dcdf609-36c9-4b22-bb16-d97e48c50f26", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("broker", "pw")
	r.ServeHTTP(rec, req)

	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), `"error": "UnknownError"`)
	assert.Contains(t, rec.Body.String(), `"description": "Could not delete service instance"`)
}
