// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomAAGUIDCreateRequestObject represents the CustomAAGUIDCreateRequestObject schema
type CustomAAGUIDCreateRequestObject struct {
	// An Authenticator Attestation Global Unique Identifier (AAGUID) is a 128-bit identifier indicating the model.
	Aaguid string `json:"aaguid,omitempty"`
	AttestationRootCertificates interface{} `json:"attestationRootCertificates,omitempty"`
	AuthenticatorCharacteristics interface{} `json:"authenticatorCharacteristics,omitempty"`
}
