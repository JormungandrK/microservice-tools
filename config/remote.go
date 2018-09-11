package config

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

// LoadRemoteConfig loads a configuration from a remote location (configURL) into an object reference.
// For convenience, it also returns the loaded object.
func LoadRemoteConfig(configURL string, configObj interface{}, templateData interface{}) (interface{}, error) {
	return LoadRemoteConfigWithLoader(configURL, NewHTTPDataLoader(&http.Client{}), configObj, templateData)
}

// LoadRemoteConfigWithLoader loads a configuration from a remote location (configURL) into an object reference using
// a DataLoader to fetch the data from the remote source.
// If config is template file it evaluates the template variable with templateData fields
// For convenience, it also returns the loaded object.
func LoadRemoteConfigWithLoader(configURL string, loader DataLoader, configObj interface{}, templateData interface{}) (interface{}, error) {
	data, err := loader(configURL)
	if err != nil {
		return nil, err
	}

	serviceConfig, err := parseConfig(data, templateData)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(serviceConfig, configObj)
	if err != nil {
		return nil, err
	}
	return configObj, nil
}

// LoadRemoteStdConfig loads a standard configuration (ServiceConfig struct) from a remote source.
func LoadRemoteStdConfig(configURL string, templateData interface{}) (*ServiceConfig, error) {
	return LoadRemoteStdConfigWithLoader(configURL, NewHTTPDataLoader(&http.Client{}), templateData)
}

// LoadRemoteStdConfigWithLoader loads a standard configuration (ServiceConfig struct) from a remote source using
// a DataLoader to fetch the data.
func LoadRemoteStdConfigWithLoader(configURL string, loader DataLoader, templateData interface{}) (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	if _, err := LoadRemoteConfigWithLoader(configURL, loader, cfg, templateData); err != nil {
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
		data, err := loadDataOverHTTP(fmt.Sprintf("%s/v1/kv/%s", consulURL, dataURI), client)
		if err != nil {
			return nil, err
		}
		return extractConsulValue(data)
	}
}

func extractConsulValue(data []byte) ([]byte, error) {
	value := []interface{}{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	base64encodedValue, err := extractConsulValueFromKVRecord(value)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(base64encodedValue)
}

func extractConsulValueFromKVRecord(record []interface{}) (string, error) {
	if len(record) == 0 {
		return "", fmt.Errorf("no value in record")
	}
	if consulValue, ok := record[0].(map[string]interface{}); ok {
		if actualValue, ok := consulValue["Value"]; ok {
			return actualValue.(string), nil
		}
		return "", fmt.Errorf("no value in record")
	}
	return "", fmt.Errorf("dont know what to do with record item %v", record[0])
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

func parseConfig(data []byte, templateData interface{}) ([]byte, error) {
	tmpl, err := template.New("config").Parse(string(data))
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, templateData)

	return buff.Bytes(), err
}
