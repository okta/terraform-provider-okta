// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PermissionConditions represents the PermissionConditions schema
// Conditions for further restricting a permission. See [Permission conditions](https://help.okta.com/okta_help.htm?type=oie&id=ext-permission-conditions).
type PermissionConditions struct {
	// Exclude attributes with specific values for the permission
	Exclude map[string]interface{} `json:"exclude,omitempty"`
	// Include attributes with specific values for the permission
	Include map[string]interface{} `json:"include,omitempty"`
}
