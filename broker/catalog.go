package broker

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/JamesClonk/elephantsql-broker/log"
	yaml "gopkg.in/yaml.v2"
)

type ServiceCatalog struct {
	Services []Service `json:"services" yaml:"services"`
}
type Service struct {
	ID                   string   `json:"id" yaml:"id"`
	Name                 string   `json:"name" yaml:"name"`
	Description          string   `json:"description" yaml:"description"`
	Bindable             bool     `json:"bindable" yaml:"bindable"`
	InstancesRetrievable bool     `json:"instances_retrievable"`
	BindingsRetrievable  bool     `json:"bindings_retrievable"`
	Tags                 []string `json:"tags" yaml:"tags"`
	Metadata             struct {
		DisplayName         string `json:"displayName" yaml:"displayName"`
		ImageURL            string `json:"imageUrl" yaml:"imageUrl"`
		LongDescription     string `json:"longDescription" yaml:"longDescription"`
		ProviderDisplayName string `json:"providerDisplayName" yaml:"providerDisplayName"`
		DocumentationURL    string `json:"documentationUrl" yaml:"documentationUrl"`
		SupportURL          string `json:"supportUrl" yaml:"supportUrl"`
	} `json:"metadata" yaml:"metadata"`
	Plans []ServicePlan `json:"plans" yaml:"plans"`
}
type ServicePlan struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Free        bool   `json:"free" yaml:"free"`
	Bindable    bool   `json:"bindable" yaml:"bindable"`
	Metadata    struct {
		DisplayName string `json:"displayName" yaml:"displayName"`
		ImageURL    string `json:"imageUrl" yaml:"imageUrl"`
		Costs       []struct {
			Amount struct {
				USDollar float64 `json:"usd" yaml:"usd"`
			} `json:"amount" yaml:"amount"`
			Unit string `json:"unit" yaml:"unit"`
		} `json:"costs" yaml:"costs"`
		Bullets          []string `json:"bullets" yaml:"bullets"`
		DedicatedService bool     `json:"dedicatedService" yaml:"dedicatedService"`
		HighAvailability bool     `json:"highAvailability" yaml:"highAvailability"`
	} `json:"metadata" yaml:"metadata"`
}

func LoadServiceCatalog(filename string) *ServiceCatalog {
	var catalog ServiceCatalog

	if _, err := os.Stat(filename); err == nil {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Errorf("could not load %s", filename)
			log.Fatalln(err)
		}
		if err := yaml.Unmarshal(data, &catalog); err != nil {
			log.Errorf("could not parse %s", filename)
			log.Fatalln(err)
		}

		// expect & hardcode certain default values
		if len(catalog.Services) < 1 {
			log.Errorln("invalid catalog, more than one service offering defined")
			log.Fatalln(catalog)
		}
		if len(catalog.Services[0].ID) == 0 {
			catalog.Services[0].ID = "elephantsql-id"
		}
		if len(catalog.Services[0].Name) == 0 {
			catalog.Services[0].Name = "elephantsql"
		}
		// displayName
		if len(catalog.Services[0].Metadata.DisplayName) == 0 {
			catalog.Services[0].Metadata.DisplayName = "ElephantSQL"
		}
		catalog.Services[0].Bindable = true
		catalog.Services[0].InstancesRetrievable = true
		catalog.Services[0].BindingsRetrievable = true

		if len(catalog.Services[0].Plans) < 1 {
			log.Errorln("invalid catalog, at least one service plan has to be defined")
			log.Fatalln(catalog)
		}

	} else {
		log.Errorf("could not load %s", filename)
		log.Fatalln(err)
	}
	return &catalog
}

func (b *Broker) Catalog(rw http.ResponseWriter, req *http.Request) {
	if b.API.DefaultRegionPlansOnly {
		region, err := b.Client.GetRegion(b.API.DefaultRegion)
		if err != nil {
			log.Errorf("could not filter plans for default region [%s]: %v", b.API.DefaultRegion, err)
		} else {
			// create new catalog without shared plans if necessary
			if !region.SharedPlans {
				newServices := make([]Service, 0)
				for _, service := range b.ServiceCatalog.Services {
					newPlans := make([]ServicePlan, 0)
					for _, plan := range service.Plans {
						if plan.Metadata.DedicatedService {
							newPlans = append(newPlans, plan)
						}
					}
					service.Plans = newPlans
					newServices = append(newServices, service)
				}
				b.ServiceCatalog.Services = newServices
			}
		}
	}
	b.write(rw, req, 200, b.ServiceCatalog)
}
