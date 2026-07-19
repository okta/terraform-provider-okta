// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicyNetworkCondition represents the PolicyNetworkCondition schema
// Specifies a network selection mode and a set of network zones to be included or excluded. If the connection parameter's data type is `ZONE`, one of the `include` or `exclude` arrays is required. Sp...
type PolicyNetworkCondition struct {
	// The zones to exclude. Required only if connection data type is `ZONE`
	Exclude []string `json:"exclude,omitempty"`
	// The zones to include. Required only if connection data type is `ZONE`
	Include []string `json:"include,omitempty"`
	Connection interface{} `json:"connection,omitempty"`
}
