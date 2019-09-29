package broker

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/JamesClonk/elephantsql-broker/api"
	"github.com/JamesClonk/elephantsql-broker/log"
	"github.com/gorilla/mux"
)

type ServiceBinding struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}
type ServiceBindingResponse struct {
	Credentials ServiceBindingResponseCredentials `json:"credentials"`
	Endpoints   []ServiceBindingResponseEndpoint  `json:"endpoints"`
	Parameters  ServiceBindingResponseParameters  `json:"parameters"`
}
type ServiceBindingResponseCredentials struct {
	URI         string `json:"uri"`
	URL         string `json:"url"`
	DatabaseURI string `json:"database_uri"`
	APIKey      string `json:"apikey"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Database    string `json:"database"`
	Scheme      string `json:"scheme"`
	Host        string `json:"host"`
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"`
}
type ServiceBindingResponseEndpoint struct {
	Host  string   `json:"host"`
	Ports []string `json:"ports"`
}
type ServiceBindingResponseParameters struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Plan   string `json:"plan"`
	Region string `json:"region"`
}

func (b *Broker) Bind(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetInstance(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorln(err)
		b.Error(rw, req, 400, "MalformedRequest", "The service instance does not exist")
		return
	}

	// verify service binding body
	if req.Body == nil {
		log.Errorf("error reading binding request: %v", req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read binding request")
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorln(err)
		log.Errorf("error reading binding request: %v", req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read binding request")
		return
	}
	if len(body) == 0 {
		body = []byte("{}")
	}

	var binding ServiceBinding
	if err := json.Unmarshal([]byte(body), &binding); err != nil {
		log.Errorln(err)
		log.Errorf("could not unmarshal binding request body: %v", string(body))
		b.Error(rw, req, 400, "MalformedRequest", "Could not unmarshal binding request")
		return
	}

	// find service_id and plan_id
	var serviceID, planID string
	for _, service := range b.ServiceCatalog.Services {
		if binding.ServiceID == service.ID {
			for _, plan := range service.Plans {
				if binding.PlanID == plan.ID {
					serviceID = service.ID
					planID = plan.ID
				}
			}
		}
	}
	if len(serviceID) == 0 {
		log.Errorf("could not find service_id [%s] for service binding creation", binding.ServiceID)
		b.Error(rw, req, 400, "MalformedRequest", "Unknown service_id")
		return
	}
	if len(planID) == 0 {
		log.Errorf("could not find plan_id [%s] for service binding creation", binding.PlanID)
		b.Error(rw, req, 400, "MalformedRequest", "Unknown plan_id")
		return
	}

	b.write(rw, req, 200, ParseBinding(instance))
}

func (b *Broker) FetchBinding(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetInstance(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorln(err)
		b.Error(rw, req, 404, "ServiceInstanceNotFound", "The service instance does not exist")
		return
	}
	b.write(rw, req, 200, ParseBinding(instance))
}

func (b *Broker) Unbind(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetInstance(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorln(err)
		b.Error(rw, req, 400, "MalformedRequest", "The service instance does not exist")
		return
	}
	b.write(rw, req, 200, map[string]string{})
}

func ParseBinding(instance *api.InstanceInfo) ServiceBindingResponse {
	credentials := ServiceBindingResponseCredentials{
		URI:         instance.URL,
		URL:         instance.URL,
		DatabaseURI: instance.URL,
		APIKey:      instance.APIKey,
	}
	endpoints := make([]ServiceBindingResponseEndpoint, 0)

	if u, err := url.Parse(instance.URL); err == nil {
		credentials.Username = u.User.Username()
		password, _ := u.User.Password()
		credentials.Password = password

		credentials.Scheme = u.Scheme
		credentials.Host = u.Host
		hostname, port, _ := net.SplitHostPort(u.Host)
		credentials.Hostname = hostname
		p, _ := strconv.Atoi(port)
		credentials.Port = p

		credentials.Database = strings.TrimPrefix(u.Path, "/")

		endpoints = append(endpoints, ServiceBindingResponseEndpoint{
			Host:  hostname,
			Ports: []string{port},
		})
	}

	return ServiceBindingResponse{
		Credentials: credentials,
		Endpoints:   endpoints,
		Parameters: ServiceBindingResponseParameters{
			ID:     instance.ID,
			Name:   instance.Name,
			Plan:   instance.Plan,
			Region: instance.Region,
		},
	}
}
