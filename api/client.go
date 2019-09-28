package api

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/JamesClonk/elephantsql-broker/config"
	"github.com/JamesClonk/elephantsql-broker/log"
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

func (c *Client) Do(req *http.Request) (int, []byte, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	log.Debugf("ElephantSQL API request [%v:%v]", req.Method, req.URL.RequestURI())

	if len(c.API.Key) > 0 {
		req.SetBasicAuth("", c.API.Key)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return 500, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, body, nil
}
