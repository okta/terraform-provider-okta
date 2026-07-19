// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ApiService represents the ApiService schema
type ApiService struct {
	// Specifies the authentication method used by the API service application when requesting an access token<br> **Note:** Only the `client_secret_basic` method is supported
	AuthenticationMethod string `json:"authenticationMethod,omitempty"`
	// A list of Okta OAuth 2.0 scopes required for the API service app to function
	Scopes []interface{} `json:"scopes,omitempty"`
	// The URL for the API service integration configuration document
	SetupInstructionsUri string `json:"setupInstructionsUri,omitempty"`
}
