package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// LoadRemoteConfig loads a configuration from a remote location (configURL) into an object reference.
// For convenience, it also returns the loaded object.
func LoadRemoteConfig(configURL string, configObj interface{}) (interface{}, error) {
	return LoadRemoteConfigWithLoader(configURL, NewHTTPDataLoader(&http.Client{}), configObj)
}

// LoadRemoteConfigWithLoader loads a configuration from a remote location (configURL) into an object reference using
// a DataLoader to fetch the data from the remote source.
// For convenience, it also returns the loaded object.
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

// LoadRemoteStdConfig loads a standard configuration (ServiceConfig struct) from a remote source.
func LoadRemoteStdConfig(configURL string) (*ServiceConfig, error) {
	return LoadRemoteStdConfigWithLoader(configURL, NewHTTPDataLoader(&http.Client{}))
}

// LoadRemoteStdConfigWithLoader loads a standard configuration (ServiceConfig struct) from a remote source using
// a DataLoader to fetch the data.
func LoadRemoteStdConfigWithLoader(configURL string, loader DataLoader) (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	if _, err := LoadRemoteConfigWithLoader(configURL, loader, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// DataLoader loads data from a remote source.
// It just specifies a contractual intrface for fetching data -
// the means of fetching data are left completely to the implementors.
type DataLoader func(dataURL string) ([]byte, error)

// NewHTTPDataLoader creates a DataLoader that fetches data from an
// HTTP server using the provided http.Client.
// The dataURL must be a full URL to the remote data.
func NewHTTPDataLoader(client *http.Client) DataLoader {
	return func(dataURL string) ([]byte, error) {
		return loadDataOverHTTP(dataURL, client)
	}
}

// NewConsulKVDataLoader creates a DataLoader that loads data from
// Consul Key-Value store.
// You must provide a URL to the Consul server and an http.Client.
// The dataURI for the data is the key under which the remote data is
// stored on the consul server.
func NewConsulKVDataLoader(consulURL string, client *http.Client) DataLoader {
	return func(dataURI string) ([]byte, error) {
		return loadDataOverHTTP(fmt.Sprintf("%s/kv/%s", consulURL, dataURI), client)
	}
}

func loadDataOverHTTP(dataURL string, client *http.Client) ([]byte, error) {
	resp, err := client.Get(dataURL)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("no response")
	}
	return ioutil.ReadAll(resp.Body)
}