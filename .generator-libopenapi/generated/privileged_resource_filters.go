// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PrivilegedResourceFilters represents the PrivilegedResourceFilters schema
type PrivilegedResourceFilters struct {
	// Array of app groups whose members might be privileged app users
	AppGroups []interface{} `json:"appGroups,omitempty"`
	// Array of organizational units where privileged app users are present
	OrganizationalUnits []interface{} `json:"organizationalUnits,omitempty"`
}
