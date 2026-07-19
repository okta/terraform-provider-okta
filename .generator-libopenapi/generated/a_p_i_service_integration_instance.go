// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// APIServiceIntegrationInstance represents the APIServiceIntegrationInstance schema
type APIServiceIntegrationInstance struct {
	Links interface{} `json:"_links,omitempty"`
	// The list of Okta management scopes granted to the API Service Integration instance. See [Okta management OAuth 2.0 scopes](/oauth2/#okta-admin-management).
	GrantedScopes []string `json:"grantedScopes,omitempty"`
	// The ID of the API Service Integration instance
	ID string `json:"id,omitempty"`
	// The name of the API service integration that corresponds with the `type` property. This is the full name of the API service integration listed in the Okta Integration Network (OIN) catalog.
	Name string `json:"name,omitempty"`
	// The URL to the API service integration configuration guide
	ConfigGuideUrl string `json:"configGuideUrl,omitempty"`
	// Timestamp when the API Service Integration instance was created
	CreatedAt string `json:"createdAt,omitempty"`
	// The user ID of the API Service Integration instance creator
	CreatedBy string `json:"createdBy,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
	// The type of the API service integration. This string is an underscore-concatenated, lowercased API service integration name. For example, `my_api_log_integration`.
	Type string `json:"type,omitempty"`
}
