// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlAssertionEncryption represents the SamlAssertionEncryption schema
// Determines if the app supports encrypted assertions
type SamlAssertionEncryption struct {
	// The key transport algorithm used to encrypt the SAML assertion
	KeyTransportAlgorithm string `json:"keyTransportAlgorithm,omitempty"`
	// A list that contains exactly one x509 encoded certificate which Okta uses to encrypt the SAML assertion with
	X5c []string `json:"x5c,omitempty"`
	// Indicates whether Okta encrypts the assertions that it sends to the Service Provider
	Enabled bool `json:"enabled,omitempty"`
	// The encryption algorithm used to encrypt the SAML assertion
	EncryptionAlgorithm string `json:"encryptionAlgorithm,omitempty"`
}
