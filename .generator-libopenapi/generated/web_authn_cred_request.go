// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// WebAuthnCredRequest represents the WebAuthnCredRequest schema
// Credential request object for the initialized credential, along with the enrollment and key identifiers to associate with the credential
type WebAuthnCredRequest struct {
	// ID for a WebAuthn preregistration factor in Okta
	AuthenticatorEnrollmentId string `json:"authenticatorEnrollmentId,omitempty"`
	// Encrypted JWE of credential request for the fulfillment provider
	CredRequestJwe string `json:"credRequestJwe,omitempty"`
	// ID for the Okta response key-pair used to encrypt and decrypt credential requests and responses
	KeyId string `json:"keyId,omitempty"`
}
