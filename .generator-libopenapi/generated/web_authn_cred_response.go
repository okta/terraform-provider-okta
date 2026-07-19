// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// WebAuthnCredResponse represents the WebAuthnCredResponse schema
// Credential response object for enrolled credential details, along with enrollment and key identifiers to associate the credential
type WebAuthnCredResponse struct {
	// ID for a WebAuthn preregistration factor in Okta
	AuthenticatorEnrollmentId string `json:"authenticatorEnrollmentId,omitempty"`
	// Encrypted JSON Web Encryption (JWE) of the credential response from the fulfillment provider
	CredResponseJwe string `json:"credResponseJwe,omitempty"`
}
