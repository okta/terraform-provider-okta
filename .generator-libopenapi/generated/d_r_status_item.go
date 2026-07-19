// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DRStatusItem represents the DRStatusItem schema
// Status whether a domain has been failed over or not
type DRStatusItem struct {
	// Domain for your org
	Domain string `json:"domain,omitempty"`
	// Indicates if the domain has been failed over
	IsFailedOver bool `json:"isFailedOver,omitempty"`
}
