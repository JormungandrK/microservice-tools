package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/JormungandrK/microservice-tools/gateway"
)

// ServiceConfig holds the full microservice configuration:
// - Configuration for registering on the API Gateway
// - Security configuration
// - Database configuration
type ServiceConfig struct {
	// Service holds the confgiuration for connecting and registering the service with the API Gateway
	Service *gateway.MicroserviceConfig `json:"service"`
	// SecurityConfig holds the security configuration
	SecurityConfig `json:"security, omitempty"`
	// DBConfig holds the database connection configuration
	DBConfig `json:"database"`
	// GatewayURL is the URL of the API Gateway
	GatewayURL string `json:"gatewayUrl"`
	// GatewayAdminURL is the administration URL of the API Gateway. Used for purposes of registration of a
	// microservice with the API gateway.
	GatewayAdminURL string `json:"gatewayAdminUrl"`
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

	// Password is the databse user password
	Password string `json:"pass,omitempty"`

	// DatabaseName is the name of the database where the server will store the collections
	DatabaseName string `json:"database,omitempty"`

	// Collections is the list of collections which should be created for the service
	Collections map[string]CollectionInfo `json:"collections,omitempty"`

	// AWSCredentials is the full path to aws credentials file
	AWSCredentials string `json:"credentials,omitempty"`

	// AWSRegion is the AWS region
	AWSRegion string `json:'awsRegion, omitemprty'`
}

// CollectionInfo holds the information about the collections
type CollectionInfo struct {
	// Indexes are the collection indexes
	Indexes []string `json:"indexes, ommitempty"`
	// HashKey is the hash key for dynamoDB table
	HashKey string `json:"hashKey, ommitempty"`
	// RangeKey is the range key for dynamoDB table
	RangeKey string `json:"rangeKey, ommitempty"`
	// EnableTTL sets the TTL for the collection
	EnableTTL bool `json:"enableTTL, ommitempty"`
	// TTL is time to live in seconds for the collection
	TTL int `json:"TTL, ommitempty"`
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
