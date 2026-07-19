// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// APIServiceIntegrationInstanceSecret represents the APIServiceIntegrationInstanceSecret schema
type APIServiceIntegrationInstanceSecret struct {
	// Status of the API Service Integration instance Secret
	Status string `json:"status"`
	Links interface{} `json:"_links"`
	// The OAuth 2.0 client secret string. The client secret string is returned in the response of a Secret creation request. In other responses (such as list, activate, or deactivate requests), the clien...
	ClientSecret string `json:"client_secret"`
	// Timestamp when the API Service Integration instance Secret was created
	Created string `json:"created"`
	// The ID of the API Service Integration instance Secret
	ID string `json:"id"`
	// Timestamp when the API Service Integration instance Secret was updated
	LastUpdated string `json:"lastUpdated"`
	// OAuth 2.0 client secret string hash
	SecretHash string `json:"secret_hash"`
}
