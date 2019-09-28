package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JamesClonk/elephantsql-broker/log"
)

type Regions []Region
type Region struct {
	Provider    string `json:"provider"`
	Region      string `json:"region"`
	Name        string `json:"name"`
	SharedPlans bool   `json:"has_shared_plans"`
}

func (c *Client) ListRegions() (Regions, error) {
	url := fmt.Sprintf("%s/regions", c.API.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	statusCode, body, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("could not list ElephantSQL regions: %s", string(body))
	}

	var regions Regions
	if err := json.Unmarshal(body, &regions); err != nil {
		log.Errorf("could not unmarshal regions response: %#v", string(body))
		return nil, err
	}
	return regions, nil
}

func (c *Client) GetRegion(key string) (*Region, error) {
	regions, err := c.ListRegions()
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		if fmt.Sprintf("%s::%s", region.Provider, region.Region) == key {
			return &region, nil
		}
	}
	return nil, fmt.Errorf("could not find ElephantSQL region [%s]", key)
}
