package notification

// Event types
const (
	// CREATE means create new service
	CREATED = "created"
	// UPDATE means update existing service
	UPDATED = "updated"
	// DELETE means delete existing service
	DELETED = "deleted"
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
