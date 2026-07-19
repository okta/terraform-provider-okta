// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CsrMetadataSubject represents the CsrMetadataSubject schema
type CsrMetadataSubject struct {
	// Common name of the subject
	CommonName string `json:"commonName,omitempty"`
	// Country name or code
	CountryName string `json:"countryName,omitempty"`
	// Locality (city) name
	LocalityName string `json:"localityName,omitempty"`
	// Name of the smaller organization, for example, the department or the division
	OrganizationalUnitName string `json:"organizationalUnitName,omitempty"`
	// Large organization name
	OrganizationName string `json:"organizationName,omitempty"`
	// State or province name
	StateOrProvinceName string `json:"stateOrProvinceName,omitempty"`
}
