package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func LoadRemoteConfig(configURL string, configObj interface{}) (interface{}, error) {
	return LoadRemoteConfigWithLoader(configURL, NewHttpDataLoader(&http.Client{}), configObj)
}

func LoadRemoteConfigWithLoader(configURL string, loader DataLoader, configObj interface{}) (interface{}, error) {
	data, err := loader(configURL)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, configObj)
	if err != nil {
		return nil, err
	}
	return configObj, nil
}

func LoadRemoteStdConfig(configURL string) (*ServiceConfig, error) {
	return LoadRemoteStdConfigWithLoader(configURL, NewHttpDataLoader(&http.Client{}))
}

func LoadRemoteStdConfigWithLoader(configURL string, loader DataLoader) (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	if _, err := LoadRemoteConfigWithLoader(configURL, loader, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

type DataLoader func(dataURL string) ([]byte, error)

func NewHttpDataLoader(client *http.Client) DataLoader {
	return func(dataURL string) ([]byte, error) {
		return loadDataOverHttp(dataURL, client)
	}
}

func NewConsulKVDataLoader(consulURL string, client *http.Client) DataLoader {
	return func(dataURI string) ([]byte, error) {
		return loadDataOverHttp(fmt.Sprintf("%s/kv/%s", consulURL, dataURI), client)
	}
}

func loadDataOverHttp(dataURL string, client *http.Client) ([]byte, error) {
	resp, err := client.Get(dataURL)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("no response")
	}
	return ioutil.ReadAll(resp.Body)
}
