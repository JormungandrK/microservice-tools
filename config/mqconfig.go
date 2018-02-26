package config

// MQConfig holds the messaging queue configuration.
type MQConfig struct {
	// Host is the remote mq host
	Host string `json:"host"`
	// Port is the port on which the remote mq server listens
	Port string `json:"port"`
	// Username to access the mq server
	Username string `json:"username"`
	// Port to access the mq server
	Password string `json:"password"`
}
