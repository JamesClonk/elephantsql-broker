package api

import (
	"io/ioutil"
	"net/http"

	"github.com/JamesClonk/elephantsql-broker/log"
)

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
