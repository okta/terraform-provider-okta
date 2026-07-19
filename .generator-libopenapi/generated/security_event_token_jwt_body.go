// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventTokenJwtBody represents the SecurityEventTokenJwtBody schema
// JSON web token body payload for a security event token sent by the SSF transmitter. For examples and more information, see [SSF Transmitter SET payload structures](https://developer.okta.com/docs/r...
type SecurityEventTokenJwtBody struct {
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
