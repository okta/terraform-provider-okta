// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventTokenJwtHeader represents the SecurityEventTokenJwtHeader schema
// JSON web token header for a security event token sent by the SSF transmitter
type SecurityEventTokenJwtHeader struct {
	// The type of content being signed or encrypted
	Typ string `json:"typ"`
	// Algorithm used to sign or encrypt the JWT
	Alg string `json:"alg"`
	// Key ID used to sign or encrypt the JWT
	Kid string `json:"kid"`
}
