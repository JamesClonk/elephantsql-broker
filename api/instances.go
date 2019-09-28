package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/JamesClonk/elephantsql-broker/log"
)

type Instances []Instance
type Instance struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Plan   string `json:"plan"`
	Region string `json:"region"`
	URL    string `json:"url"`
}

type InstanceInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Plan   string `json:"plan"`
	Region string `json:"region"`
	URL    string `json:"url"`
	APIKey string `json:"apikey"`
	Ready  bool   `json:"ready"`
}

type CreateInstanceResponse struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	APIKey string `json:"apikey"`
}

func (c *Client) CreateInstance(name, plan, region string) (*CreateInstanceResponse, error) {
	if len(region) == 0 {
		region = c.API.DefaultRegion
	}
	values := url.Values{
		"name":   {name},
		"plan":   {plan},
		"region": {region},
	}

	url := fmt.Sprintf("%s/instances", c.API.URL)
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	statusCode, body, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusCreated {
		return nil, fmt.Errorf("could not create ElephantSQL instance [%s]: %s", name, string(body))
	}

	var response CreateInstanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Errorf("could not unmarshal create instance response: %#v", string(body))
		return nil, err
	}
	return &response, nil
}

func (c *Client) DeleteInstance(instance InstanceInfo) error {
	url := fmt.Sprintf("%s/instances/%d", c.API.URL, instance.ID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	statusCode, body, err := c.Do(req)
	if err != nil {
		return err
	}

	if statusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete ElephantSQL instance [%s]: %s", instance.Name, string(body))
	}
	return nil
}

func (c *Client) GetInstanceInfo(id int) (*InstanceInfo, error) {
	url := fmt.Sprintf("%s/instances/%d", c.API.URL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	statusCode, body, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get ElephantSQL instance [%d] info: %s", id, string(body))
	}

	var info InstanceInfo
	if err := json.Unmarshal(body, &info); err != nil {
		log.Errorf("could not unmarshal instance info response: %#v", string(body))
		return nil, err
	}
	return &info, nil
}

func (c *Client) ListInstances() (Instances, error) {
	url := fmt.Sprintf("%s/instances", c.API.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	statusCode, body, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("could not list ElephantSQL instances: %s", string(body))
	}

	var instances Instances
	if err := json.Unmarshal(body, &instances); err != nil {
		log.Errorf("could not unmarshal instances response: %#v", string(body))
		return nil, err
	}
	return instances, nil
}

func (c *Client) GetInstance(name string) (*InstanceInfo, error) {
	instances, err := c.ListInstances()
	if err != nil {
		return nil, err
	}

	for _, instance := range instances {
		if instance.Name == name {
			return c.GetInstanceInfo(instance.ID)
		}
	}
	return nil, fmt.Errorf("could not find ElephantSQL instance by name [%s]", name)
}
