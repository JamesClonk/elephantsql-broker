package broker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swisscom/backman/log"
)

type ServiceInstanceProvisioning struct {
	ServiceID  string `json:"service_id"`
	PlanID     string `json:"plan_id"`
	Parameters struct {
		Region string `json:"region"`
	} `json:"parameters"`
}
type ServiceInstanceProvisioningResponse struct {
	DashboardURL string `json:"dashboard_url"`
}

type ServiceInstanceFetchResponse struct {
	ServiceID    string                                 `json:"service_id,omitempty"`
	PlanID       string                                 `json:"plan_id,omitempty"`
	DashboardURL string                                 `json:"dashboard_url"`
	Parameters   ServiceInstanceFetchResponseParameters `json:"parameters"`
}
type ServiceInstanceFetchResponseParameters struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Plan   string `json:"plan"`
	Region string `json:"region"`
}

func (b *Broker) ProvisionInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorln(err)
		log.Errorf("error reading provisioning request: %v", req)
		b.Error(rw, req, 400, "MalformedRequest", "Could not read provisioning request")
		return
	}
	if len(body) == 0 {
		body = []byte("{}")
	}

	var provisioning ServiceInstanceProvisioning
	if err := json.Unmarshal([]byte(body), &provisioning); err != nil {
		log.Errorln(err)
		log.Errorf("could not unmarshal provisioning request body: %v", string(body))
		b.Error(rw, req, 400, "MalformedRequest", "Could not unmarshal provisioning request")
		return
	}

	// find plan name
	var instancePlan string
	for _, service := range b.ServiceCatalog.Services {
		if provisioning.ServiceID == service.ID {
			for _, plan := range service.Plans {
				if provisioning.PlanID == plan.ID {
					instancePlan = plan.Name
				}
			}
		}
	}
	if len(instancePlan) == 0 {
		log.Errorf("could not find plan_id [%s] for service instance creation", provisioning.PlanID)
		b.Error(rw, req, 400, "MalformedRequest", "Unknown plan_id")
		return
	}

	// use default region if parameter was not provided
	instanceRegion := provisioning.Parameters.Region
	if len(instanceRegion) == 0 {
		instanceRegion = b.API.DefaultRegion
	}

	// check if conflicting service instance already exists
	instance, err := b.Client.GetInstance(instanceID)
	if err == nil {
		if instance.Name != instanceID || instance.Plan != instancePlan || instance.Region != instanceRegion {
			log.Errorf("service instance [%s] already exists with conflicting attributes", instanceID)
			b.Error(rw, req, 409, "Conflict", "The service instance already exists with different attributes")
			return
		}
		if instance.Name == instanceID && instance.Plan == instancePlan && instance.Region == instanceRegion {
			log.Debugf("service instance [%s] already exists", instanceID)
			// response JSON
			provisioningResponse := ServiceInstanceProvisioningResponse{
				DashboardURL: fmt.Sprintf("https://customer.elephantsql.com/instance/%d/sso", instance.ID),
			}
			b.write(rw, req, 200, provisioningResponse)
			return
		}
	}

	// create
	resultingInstance, err := b.Client.CreateInstance(instanceID, instancePlan, instanceRegion)
	if err != nil {
		log.Errorln(err)
		b.Error(rw, req, 500, "UnknownError", "Could not create service instance")
		return
	}

	// response JSON
	provisioningResponse := ServiceInstanceProvisioningResponse{
		DashboardURL: fmt.Sprintf("https://customer.elephantsql.com/instance/%d/sso", resultingInstance.ID),
	}
	b.write(rw, req, 201, provisioningResponse)
}

func (b *Broker) FetchInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetInstance(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorln(err)
		b.Error(rw, req, 404, "ServiceInstanceNotFound", "The service instance does not exist")
		return
	}

	// find service_id and plan_id
	var serviceID, planID string
	for _, service := range b.ServiceCatalog.Services {
		for _, plan := range service.Plans {
			if plan.Name == instance.Plan {
				serviceID = service.ID
				planID = plan.ID
			}
		}
	}

	// response JSON
	fetchResponse := ServiceInstanceFetchResponse{
		ServiceID:    serviceID,
		PlanID:       planID,
		DashboardURL: fmt.Sprintf("https://customer.elephantsql.com/instance/%d/sso", instance.ID),
		Parameters: ServiceInstanceFetchResponseParameters{
			ID:     instance.ID,
			Name:   instance.Name,
			Plan:   instance.Plan,
			Region: instance.Region,
		},
	}
	b.write(rw, req, 200, fetchResponse)
}

func (b *Broker) DeprovisionInstance(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instanceID"]

	instance, err := b.Client.GetInstance(instanceID)
	if err != nil || instance.Name != instanceID {
		log.Errorln(err)
		b.Error(rw, req, 410, "MissingServiceInstance", "The service instance does not exist")
		return
	}

	if err := b.Client.DeleteInstance(*instance); err != nil {
		log.Errorln(err)
		b.Error(rw, req, 500, "UnknownError", "Could not delete service instance")
		return
	}
	b.write(rw, req, 200, map[string]string{})
}
