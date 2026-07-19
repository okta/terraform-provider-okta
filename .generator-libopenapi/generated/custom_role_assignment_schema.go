// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomRoleAssignmentSchema represents the CustomRoleAssignmentSchema schema
type CustomRoleAssignmentSchema struct {
	// Resource set ID
	Resource-set string `json:"resource-set"`
	// Custom role ID
	Role string `json:"role"`
	// Specify a [standard admin role](/openapi/okta-management/guides/roles/#standard-roles), an [IAM-based standard role](/openapi/okta-management/guides/roles/#iam-based-standard-roles), or `CUSTOM` fo...
	Type string `json:"type"`
}
