package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Microkubes/microservice-tools/gateway"
)

// ServiceConfig holds the full microservice configuration:
// - Configuration for registering on the API Gateway
// - Security configuration
// - Database configuration
type ServiceConfig struct {
	// Service holds the confgiuration for connecting and registering the service with the API Gateway
	Service *gateway.MicroserviceConfig `json:"service"`
	// SecurityConfig holds the security configuration
	SecurityConfig `json:"security,omitempty"`
	// DBConfig holds the database connection configuration
	DBConfig `json:"database"`
	// GatewayURL is the URL of the API Gateway
	GatewayURL string `json:"gatewayUrl"`
	// GatewayAdminURL is the administration URL of the API Gateway. Used for purposes of registration of a
	// microservice with the API gateway.
	GatewayAdminURL string `json:"gatewayAdminUrl"`
	// ContainerManager is the platform for managing containerized services
	// Can be swarm or kubernetes
	ContainerManager string `json:"containerManager,omitempty"`
}

// DBConfig holds the database configuration parameters.
type DBConfig struct {
	// DBname is the database name (mongodb/dynamodb)
	DBName string `json:"dbName"`

	// DB Info holds the database connection configuration
	DBInfo `json:"dbInfo"`
}

// DBInfo holds the database connection configuration
type DBInfo struct {
	// Host is the database host+port URL
	Host string `json:"host,omitempty"`

	// Username is the username used to access the database
	Username string `json:"user,omitempty"`

	// Password is the database user password
	Password string `json:"pass,omitempty"`

	// DatabaseName is the name of the database where the server will store the collections
	DatabaseName string `json:"database,omitempty"`

	// AWSCredentials is the full path to aws credentials file
	AWSCredentials string `json:"credentials,omitempty"`

	// AWSEndpoint is the full path to aws credentials file
	AWSEndpoint string `json:"endpoint,omitempty"`

	// AWSRegion is the AWS region
	AWSRegion string `json:"awsRegion,omitempty"`
}

// LoadConfig loads the service configuration from a file.
func LoadConfig(confFile string) (*ServiceConfig, error) {
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	conf := &ServiceConfig{}
	err = json.Unmarshal(data, conf)
	return conf, err
}

// LoadConfigAs loads a generic configuration from a JSON file into predefined structure.
func LoadConfigAs(confFile string, conf interface{}) error {
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, conf)
}

func readFileAndMerge(confFile string, variables interface{}) ([]byte, error) {
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	if variables == nil {
		return data, nil
	}
	return parseConfig(data, variables)
}

// LoadConfigAsTypeAndMerge loads configuration template from a file and merges with the provided variables.
// The configuration is unmarshalled into the provided configuration interface.
func LoadConfigAsTypeAndMerge(confFile string, conf interface{}, variables interface{}) error {
	data, err := readFileAndMerge(confFile, variables)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, conf)
}

// LoadConfigAndMerge loads configuration template from a file and merges the template with the provided variables.
// The merged configuration is then unmarshalled into a standard ServiceConfig structure.
func LoadConfigAndMerge(confFile string, variables interface{}) (*ServiceConfig, error) {
	data, err := readFileAndMerge(confFile, variables)
	if err != nil {
		return nil, err
	}
	conf := &ServiceConfig{}
	err = json.Unmarshal(data, conf)
	return conf, err
}
