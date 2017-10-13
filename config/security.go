package config

type SecurityConfig struct {
	KeysDir      string `json:"keysDir"`
	JWTConfig    `json:"jwt"`
	SAMLConfig   `json:"saml"`
	OAuth2Config `json:"oauth2"`
}

type JWTConfig struct {
	Name        string
	Description string
	TokenURL    string `json:"tokenUrl"`
}

type SAMLConfig struct {
	CertFile               string `json:"certFile"`
	KeyFile                string `json:"keyFile"`
	IdentityProviderURL    string `json:"identityProviderUrl"`
	UserServiceURL         string `json:"userServiceUrl"`
	RegistrationServiceURL string `json:"registrationServiceUrl"`
}

type OAuth2Config struct {
	TokenURL         string `json:"tokenUrl"`
	AuthorizationURL string `json:"authorizeUrl"`
	Description      string `json:"description"`
}
