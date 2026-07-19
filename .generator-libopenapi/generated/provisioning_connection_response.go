// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProvisioningConnectionResponse represents the ProvisioningConnectionResponse schema
type ProvisioningConnectionResponse struct {
	// Base URL
	BaseUrl string `json:"baseUrl,omitempty"`
	Profile interface{} `json:"profile"`
	Status interface{} `json:"status"`
	Links interface{} `json:"_links,omitempty"`
	AuthScheme interface{} `json:"authScheme,omitempty"`
}
