// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DomainListResponse represents the DomainListResponse schema
// Defines a list of domains with a subset of the properties for each domain
type DomainListResponse struct {
	// Each element of the array defines an individual domain
	Domains []interface{} `json:"domains,omitempty"`
}
