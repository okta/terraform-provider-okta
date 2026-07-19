// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogIpServiceCategory represents the LogIpServiceCategory schema
// Describes the IP service category associated with an IP.
type LogIpServiceCategory struct {
	// The type of service provided from this IP address
	Type string `json:"type,omitempty"`
	// Indicates whether the service is an anonymizer
	IsAnonymous bool `json:"isAnonymous,omitempty"`
	// The name of the associated operator
	Operator string `json:"operator,omitempty"`
}
