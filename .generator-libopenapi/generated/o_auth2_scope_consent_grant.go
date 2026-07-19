// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ScopeConsentGrant represents the OAuth2ScopeConsentGrant schema
// Grant object that represents an app consent scope grant
type OAuth2ScopeConsentGrant struct {
	// Client ID of the app integration
	ClientId string `json:"clientId,omitempty"`
	Created interface{} `json:"created,omitempty"`
	// ID of the Grant object
	ID string `json:"id,omitempty"`
	// The issuer of your org authorization server. This is typically your Okta domain.
	Issuer string `json:"issuer"`
	LastUpdated interface{} `json:"lastUpdated,omitempty"`
	Source interface{} `json:"source,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	CreatedBy interface{} `json:"createdBy,omitempty"`
	// The name of the [Okta scope](https://developer.okta.com/docs/api/oauth2/#oauth-20-scopes) for which consent is granted
	ScopeId string `json:"scopeId"`
	// User ID that granted consent (if `source` is `END_USER`)
	UserId string `json:"userId,omitempty"`
	// Embedded resources related to the Grant
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
}
