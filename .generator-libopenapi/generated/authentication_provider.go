// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticationProvider represents the AuthenticationProvider schema
// Specifies the authentication provider that validates the user's password credential. The user's current provider is managed by the **Delegated Authentication** settings for your org. The provider o...
type AuthenticationProvider struct {
	// The name of the authentication provider
	Name string `json:"name,omitempty"`
	Type interface{} `json:"type,omitempty"`
}
