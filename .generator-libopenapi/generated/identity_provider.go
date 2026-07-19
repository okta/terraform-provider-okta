// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdentityProvider represents the IdentityProvider schema
type IdentityProvider struct {
	IssuerMode interface{} `json:"issuerMode,omitempty"`
	LastUpdated interface{} `json:"lastUpdated,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Created interface{} `json:"created,omitempty"`
	// Unique key for the IdP
	ID string `json:"id,omitempty"`
	// Unique name for the IdP
	Name string `json:"name,omitempty"`
	Policy interface{} `json:"policy,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
	// IdP-specific protocol settings for endpoints, bindings, and algorithms used to connect with the IdP and validate messages
	Protocol interface{} `json:"protocol,omitempty"`
	Type interface{} `json:"type,omitempty"`
	Links map[string]interface{} `json:"_links,omitempty"`
}
