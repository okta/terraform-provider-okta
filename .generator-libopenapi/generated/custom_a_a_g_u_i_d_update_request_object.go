// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomAAGUIDUpdateRequestObject represents the CustomAAGUIDUpdateRequestObject schema
type CustomAAGUIDUpdateRequestObject struct {
	AuthenticatorCharacteristics interface{} `json:"authenticatorCharacteristics,omitempty"`
	// The product name associated with this AAGUID.
	Name string `json:"name,omitempty"`
	AttestationRootCertificates interface{} `json:"attestationRootCertificates,omitempty"`
}
