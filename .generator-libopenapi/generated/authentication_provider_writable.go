// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticationProviderWritable represents the AuthenticationProviderWritable schema
// Specifies the authentication provider that validates the user password credential. The user's current provider is managed by the **Delegated Authentication** settings in your org. See [Create user ...
type AuthenticationProviderWritable struct {
	// The name of the authentication provider
	Name string `json:"name,omitempty"`
	Type interface{} `json:"type,omitempty"`
}
