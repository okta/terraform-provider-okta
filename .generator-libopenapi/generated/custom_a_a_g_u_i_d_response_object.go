// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomAAGUIDResponseObject represents the CustomAAGUIDResponseObject schema
type CustomAAGUIDResponseObject struct {
	// A unique 128-bit identifier that's assigned to a specific model of security key or authenticator
	Aaguid string `json:"aaguid,omitempty"`
	AttestationRootCertificates interface{} `json:"attestationRootCertificates,omitempty"`
	AuthenticatorCharacteristics interface{} `json:"authenticatorCharacteristics,omitempty"`
	// The product name associated with the AAGUID
	Name string `json:"name,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
