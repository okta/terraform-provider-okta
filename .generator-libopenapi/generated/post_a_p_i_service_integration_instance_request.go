// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// postAPIServiceIntegrationInstanceRequest represents the postAPIServiceIntegrationInstanceRequest schema
type postAPIServiceIntegrationInstanceRequest struct {
	Properties interface{} `json:"properties,omitempty"`
	// The type of the API service integration. This string is an underscore-concatenated, lowercased API service integration name. For example, `my_api_log_integration`.
	Type string `json:"type"`
	// The list of Okta management scopes granted to the API Service Integration instance. See [Okta management OAuth 2.0 scopes](/oauth2/#okta-admin-management).
	GrantedScopes []string `json:"grantedScopes"`
}
