// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CreateResourceSetRequest represents the CreateResourceSetRequest schema
type CreateResourceSetRequest struct {
	// Description of the resource set
	Description string `json:"description"`
	// Unique name for the resource set
	Label string `json:"label"`
	// The endpoint (URL) that references all resource objects included in the resource set. Resources are identified by either an Okta Resource Name (ORN) or by a REST URL format. See [Okta Resource Name...
	Resources []string `json:"resources"`
}
