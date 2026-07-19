// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlAttributeStatementGroup represents the SamlAttributeStatementGroup schema
// `GROUP` attribute statements
type SamlAttributeStatementGroup struct {
	// The operation to filter groups based on `filterValue`
	FilterType string `json:"filterType,omitempty"`
	// Filter the groups based on a specific value.
	FilterValue string `json:"filterValue,omitempty"`
	// The name of the group attribute in your app. The attribute name must be unique across all user and group attribute statements.
	Name string `json:"name,omitempty"`
	// The name format of the group attribute. Supported values:
	Namespace string `json:"namespace,omitempty"`
	// The type of attribute statements object
	Type string `json:"type,omitempty"`
}
