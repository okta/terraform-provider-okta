// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventTokenRequestJwtHeader represents the SecurityEventTokenRequestJwtHeader schema
// JSON web token header for a security event token
type SecurityEventTokenRequestJwtHeader struct {
	// Algorithm used to sign or encrypt the JWT
	Alg string `json:"alg"`
	// Key ID used to sign or encrypt the JWT
	Kid string `json:"kid"`
	// The type of content being signed or encrypted
	Typ string `json:"typ"`
}
