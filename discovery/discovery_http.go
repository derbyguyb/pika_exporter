package discovery

import (
	"encoding/json"
	"errors"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var invalidResponse = errors.New("invalid response")

func NewHttpDiscovery(httpPath string) (*httpDiscovery, error) {
	req, err := newHttpRequest(httpPath)
	if err != nil {
		return nil, err
	}

	discover := &httpDiscovery{cli: &http.Client{
		Timeout: 10 * time.Second,
	}, req: req}

	err = discover.updateInstances()
	if err != nil {
		return nil, err
	}

	discover.startRegularUpdate()
	return discover, nil
}

func newHttpRequest(httpPath string) (*http.Request, error) {
	encodedUrl, err := httpUrlEncoded(httpPath)
	if err != nil {
		return nil, err
	}
	return http.NewRequest("GET", encodedUrl, nil)
}

func httpUrlEncoded(str string) (string, error) {
	var Url *url.URL
	Url, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	parameters := url.Values{}
	for key, values := range Url.Query() {
		for _, value := range values {
			parameters.Add(key, value)
		}
	}
	Url.RawQuery = parameters.Encode()
	return Url.String(), nil
}

type httpDiscovery struct {
	cli       *http.Client
	req       *http.Request
	instances []Instance
}

func (http *httpDiscovery) GetInstances() []Instance {
	return http.instances
}

func (http *httpDiscovery) startRegularUpdate() {
	go func() {
		tickerChan := time.Tick(time.Minute)
		for {
			e := http.updateInstances()
			if e != nil {
				log.Error("regular update:", e)
			}
			<-tickerChan
		}
	}()
}

func (http *httpDiscovery) updateInstances() error {
	rs, err := http.getResponse()
	if err != nil {
		return err
	}

	httpJson := &httpJson{}
	if err = json.Unmarshal(rs, &httpJson); err != nil {
		return err
	}

	if httpJson.PikaInstances == nil {
		return invalidResponse
	}

	http.instances = httpJson.PikaInstances
	return nil
}

func (http *httpDiscovery) getResponse() ([]byte, error) {
	resp, err := http.cli.Do(http.req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if resp != nil {
			err := resp.Body.Close()
			if err != nil {
				log.Error(err)
			}
		}
	}()

	if resp.StatusCode != 200 {
		return nil, invalidResponse
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type httpJson struct {
	PikaInstances []Instance
}
