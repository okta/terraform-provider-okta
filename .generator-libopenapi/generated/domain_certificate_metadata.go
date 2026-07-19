// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DomainCertificateMetadata represents the DomainCertificateMetadata schema
// Certificate metadata for the domain
type DomainCertificateMetadata struct {
	// Certificate fingerprint
	Fingerprint string `json:"fingerprint,omitempty"`
	// Certificate subject
	Subject string `json:"subject,omitempty"`
	// Certificate expiration
	Expiration string `json:"expiration,omitempty"`
}
