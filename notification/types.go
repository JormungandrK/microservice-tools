package notification

// Event types
const (
	// CREATE means create new service
	CREATE = "create"
	// UPDATE means update existing service
	UPDATE = "update"
	// DELETE means delete existing service
	DELETE = "delete"
)

// Object types
const (
	// Event object for microservice-organization
	ORGANIZATION = "organization"
	// Event object for schema-management
	SCHEMA = "schema"
	// Event object for microservice-resources
	RESOURCE = "resource"
	// Event object for microservice-data
	DATA = "data"
)
