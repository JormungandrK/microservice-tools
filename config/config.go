package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/JormungandrK/microservice-tools/gateway"
)

type ServiceConfig struct {
	Service        *gateway.MicroserviceConfig `json:"service"`
	SecurityConfig `json:"security, omitempty"`
	DBConfig       `json:"database"`
	GatewayURL     string `json:"gatewayUrl"`
}

// DBConfig holds the database configuration parameters.
type DBConfig struct {

	// Host is the database host+port URL
	Host string `json:"host,omitempty"`

	// Username is the username used to access the database
	Username string `json:"user,omitempty"`

	// Password is the databse user password
	Password string `json:"pass,omitempty"`

	// DatabaseName is the name of the database where the server will store the collections
	DatabaseName string `json:"database,omitempty"`
}

func LoadConfig(confFile string) (*ServiceConfig, error) {
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	conf := &ServiceConfig{}
	err = json.Unmarshal(data, conf)
	return conf, err
}
