// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CreateIamRoleRequest represents the CreateIamRoleRequest schema
type CreateIamRoleRequest struct {
	// Unique label for the role
	Label string `json:"label"`
	// Array of permissions that the role grants. See [Permissions](/openapi/okta-management/guides/permissions).
	Permissions []interface{} `json:"permissions"`
	// Description of the role
	Description string `json:"description"`
}
