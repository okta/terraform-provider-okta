// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2Scope represents the OAuth2Scope schema
type OAuth2Scope struct {
	// Name of the end user displayed in a consent dialog
	DisplayName string `json:"displayName,omitempty"`
	// Scope object ID
	ID string `json:"id,omitempty"`
	// Scope name
	Name string `json:"name"`
	Links interface{} `json:"_links,omitempty"`
	// Indicates if this Scope is a default scope
	Default bool `json:"default,omitempty"`
	MetadataPublish interface{} `json:"metadataPublish,omitempty"`
	// Indicates whether the Scope is optional. When set to `true`, the user can skip consent for the scope.
	Optional bool `json:"optional,omitempty"`
	// Indicates if Okta created the Scope
	System bool `json:"system,omitempty"`
	Consent interface{} `json:"consent,omitempty"`
	// Description of the Scope
	Description string `json:"description,omitempty"`
}
