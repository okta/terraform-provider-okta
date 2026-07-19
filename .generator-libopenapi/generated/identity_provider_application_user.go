// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdentityProviderApplicationUser represents the IdentityProviderApplicationUser schema
type IdentityProviderApplicationUser struct {
	LastUpdated interface{} `json:"lastUpdated,omitempty"`
	// IdP-specific profile for the user.  IdP user profiles are IdP-specific but may be customized by the Profile Editor in the Admin Console.  > **Note:** Okta variable names have reserved characters th...
	Profile map[string]interface{} `json:"profile,omitempty"`
	// Embedded resources related to the IdP user
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Links map[string]interface{} `json:"_links,omitempty"`
	Created interface{} `json:"created,omitempty"`
	// Unique IdP-specific identifier for the user
	ExternalId string `json:"externalId,omitempty"`
	// Unique key of the user
	ID string `json:"id,omitempty"`
}
