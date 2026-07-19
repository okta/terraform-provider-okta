// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticationMethodChain represents the AuthenticationMethodChain schema
type AuthenticationMethodChain struct {
	// The next steps of the authentication method chain. This is an array of `AuthenticationMethodChain`. Only supports one item in the array.
	Next []map[string]interface{} `json:"next,omitempty"`
	// Specifies how often the user is prompted for authentication using duration format for the time period. For example, `PT2H30M` for two and a half hours. This parameter can't be set at the same time ...
	ReauthenticateIn string `json:"reauthenticateIn,omitempty"`
	AuthenticationMethods []interface{} `json:"authenticationMethods,omitempty"`
}
