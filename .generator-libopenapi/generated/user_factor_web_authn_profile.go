// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorWebAuthnProfile represents the UserFactorWebAuthnProfile schema
type UserFactorWebAuthnProfile struct {
	// Human-readable name of the authenticator  > **Note:** This name is set from the AAGUID metadata during enrollment. It can't be changed in the Admin Console or by using any Okta APIs.
	AuthenticatorName string `json:"authenticatorName,omitempty"`
	// ID for the factor credential
	CredentialId string `json:"credentialId,omitempty"`
}
