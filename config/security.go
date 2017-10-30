package config

// SecurityConfig holds the security configuration.
// The subsections for JWT, SAML and OAuth2 are optional. If a
// subsection is ommited, then the appropriate security will not be
// configured and used by the security chain.
type SecurityConfig struct {

	// KeysDir is the loacation of the directory holding the private-public key pairs.
	KeysDir string `json:"keysDir"`

	// JWTConfig holds the JWT configuration. If ommited the JWT security will not be used.
	*JWTConfig `json:"jwt,omitempty"`

	// SAMLConfig holds the SAML configuration. If ommited the SAML security will not be used.
	*SAMLConfig `json:"saml,omitempty"`

	//OAuth2Config holds the OAuth2 configuration. If ommited the OAuth2 security will not be used.
	*OAuth2Config `json:"oauth2,omitempty"`
}

// JWTConfig holds the JWT configuration.
type JWTConfig struct {

	// Name is the name of the JWT middleware. Used in error messages.
	Name string

	// Description holds the description for the middleware. Used for documentation purposes.
	Description string

	// TokenURL is the URL of the JWT token provider. Use a full URL here.
	TokenURL string `json:"tokenUrl"`
}

// SAMLConfig holds the SAML configuration.
type SAMLConfig struct {

	// CertFile is the location of the certificate file.
	CertFile string `json:"certFile"`

	// KeyFile is the location of the key file.
	KeyFile string `json:"keyFile"`

	// IdentityProviderURL is the URL of the SAML Identity Provider server. User a full URL here.
	IdentityProviderURL string `json:"identityProviderUrl"`

	// UserServiceURL is the URL of the user microservice. This should be the public url (usually over the Gateway).
	UserServiceURL string `json:"userServiceUrl"`

	// RegistrationServiceURL is the URL of the registration service. This should be the public registration URL (usually over the Gateway).
	RegistrationServiceURL string `json:"registrationServiceUrl"`

	// RootURL is the base URL of the microservice
	RootURL string `json:"rootURL"`
}

// OAuth2Config holds the OAuth2 configuration.
type OAuth2Config struct {

	// TokenURL is the path of the token endpoint. Usually "/oauth2/token".
	TokenURL string `json:"tokenUrl"`

	// AuthorizationURL is the path of the authorize endpoint. Usually "/oauth2/authorize".
	AuthorizationURL string `json:"authorizeUrl"`

	// Description is the description of the middleware. Used for documentation purposes.
	Description string `json:"description"`
}
