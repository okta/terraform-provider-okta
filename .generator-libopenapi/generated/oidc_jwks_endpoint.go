// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OidcJwksEndpoint represents the OidcJwksEndpoint schema
// Endpoint for the JSON Web Key Set (JWKS) document. This document contains signing keys that are used to validate the signatures from the provider. For more information on JWKS, see [JSON Web Key](h...
type OidcJwksEndpoint struct {
	Binding interface{} `json:"binding,omitempty"`
	// URL of the endpoint to the JWK Set
	Url string `json:"url,omitempty"`
}
