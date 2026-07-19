// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OptInStatusResponse represents the OptInStatusResponse schema
type OptInStatusResponse struct {
	// The entitlement management opt-in status for the Admin Console
	OptInStatus string `json:"optInStatus,omitempty"`
	// Link relations available
	Links map[string]interface{} `json:"_links,omitempty"`
}
