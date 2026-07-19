// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ThreatInsightConfiguration represents the ThreatInsightConfiguration schema
type ThreatInsightConfiguration struct {
	// Timestamp when the ThreatInsight Configuration object was created
	Created *time.Time `json:"created,omitempty"`
	// Accepts a list of [Network Zone](/openapi/okta-management/management/networkzone/) IDs. IPs in the excluded network zones aren't logged or blocked. This ensures that traffic from known, trusted IPs...
	ExcludeZones []string `json:"excludeZones,omitempty"`
	// Timestamp when the ThreatInsight Configuration object was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Specifies how Okta responds to authentication requests from suspicious IP addresses
	Action string `json:"action"`
}
