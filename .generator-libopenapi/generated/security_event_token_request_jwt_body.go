// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventTokenRequestJwtBody represents the SecurityEventTokenRequestJwtBody schema
// JSON web token body payload for a security event token
type SecurityEventTokenRequestJwtBody struct {
	// Audience
	Aud string `json:"aud"`
	Events interface{} `json:"events"`
	// Token issue time (UNIX timestamp)
	Iat int64 `json:"iat"`
	// Token issuer
	Iss string `json:"iss"`
	// Token ID
	Jti string `json:"jti"`
}
