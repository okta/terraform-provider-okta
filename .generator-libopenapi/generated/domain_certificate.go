// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DomainCertificate represents the DomainCertificate schema
// Defines the properties of the certificate
type DomainCertificate struct {
	Type interface{} `json:"type"`
	// Certificate content
	Certificate string `json:"certificate"`
	// Certificate chain
	CertificateChain string `json:"certificateChain"`
	// Certificate private key
	PrivateKey string `json:"privateKey"`
}
