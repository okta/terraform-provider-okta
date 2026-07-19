// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlSpCertificate represents the SamlSpCertificate schema
// The certificate that Okta uses to validate Single Logout (SLO) requests and responses
type SamlSpCertificate struct {
	// A list that contains exactly one x509 encoded certificate
	X5c []string `json:"x5c,omitempty"`
}
