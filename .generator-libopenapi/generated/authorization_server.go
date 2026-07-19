// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// AuthorizationServer represents the AuthorizationServer schema
type AuthorizationServer struct {
	// Indicates which value is specified in the issuer of the tokens that a custom authorization server returns: the Okta org domain URL or a custom domain URL.  `issuerMode` is visible if you have a cus...
	IssuerMode string `json:"issuerMode,omitempty"`
	// The name of the custom authorization server
	Name string `json:"name,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// The recipients that the tokens are intended for. This becomes the `aud` claim in an access token. Okta currently supports only one audience.
	Audiences []string `json:"audiences,omitempty"`
	Created *time.Time `json:"created,omitempty"`
	// The ID of the custom authorization server
	ID string `json:"id,omitempty"`
	// The complete URL for the custom authorization server. This becomes the `iss` claim in an access token.
	Issuer string `json:"issuer,omitempty"`
	Jwks interface{} `json:"jwks,omitempty"`
	// URL string that references a JSON Web Key Set for encrypting JWTs minted by the custom authorization server
	JwksUri string `json:"jwks_uri,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	AccessTokenEncryptedResponseAlgorithm interface{} `json:"accessTokenEncryptedResponseAlgorithm,omitempty"`
	Credentials interface{} `json:"credentials,omitempty"`
	// The description of the custom authorization server
	Description string `json:"description,omitempty"`
}
