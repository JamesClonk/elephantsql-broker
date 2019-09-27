package api

import (
	"net/http"
	"sync"

	"github.com/JamesClonk/elephantsql-broker/config"
	"github.com/JamesClonk/elephantsql-broker/util"
)

type Client struct {
	API    config.API
	Mutex  *sync.Mutex
	Client *http.Client
}

func NewClient(c *config.Config) *Client {
	client := &Client{
		API:    c.API,
		Mutex:  &sync.Mutex{},
		Client: util.NewHttpClient(c),
	}
	return client
}
