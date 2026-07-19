// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdPCsr represents the IdPCsr schema
// Defines a CSR for a signature or decryption credential for an IdP
type IdPCsr struct {
	Created interface{} `json:"created,omitempty"`
	// Base64-encoded CSR in DER format
	Csr string `json:"csr,omitempty"`
	// Unique identifier for the CSR
	ID string `json:"id,omitempty"`
	// Cryptographic algorithm family for the CSR's keypair
	Kty string `json:"kty,omitempty"`
	Links map[string]interface{} `json:"_links,omitempty"`
}
