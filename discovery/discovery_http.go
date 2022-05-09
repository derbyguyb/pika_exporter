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

var _IllegalResponse = errors.New("Illegal response")

func NewHttpDiscovery(httpPath string) (*httpDiscovery, error) {
	encodedUrl, err := httpUrlEncoded(httpPath)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", encodedUrl, nil)
	if err != nil {
		return nil, err
	}

	discover := &httpDiscovery{cli: &http.Client{
		Timeout: 10 * time.Second,
	}, req: req}
	err = discover.init()
	if err != nil {
		return nil, err
	}

	return discover, nil
}

type httpDiscovery struct {
	cli       *http.Client
	req       *http.Request
	instances []Instance
}

func (http *httpDiscovery) GetInstances() []Instance {
	go http.init()
	return http.instances
}

func (http *httpDiscovery) init() error {
	rs, err := http.getResponse()
	if err != nil {
		log.Error(err)
		return err
	}

	httpJson := &httpJson{}
	if err = json.Unmarshal(rs, &httpJson); err != nil {
		log.Error(err)
		return err
	}

	if httpJson == nil || httpJson.PikaInstances == nil {
		log.Error(_IllegalResponse)
		return _IllegalResponse
	}

	http.instances = httpJson.PikaInstances
	return nil
}

func (http *httpDiscovery) getResponse() ([]byte, error) {
	resp, err := http.cli.Do(http.req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
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

type httpJson struct {
	PikaInstances []Instance
}
