// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DomainResponse represents the DomainResponse schema
// The properties that define an individual domain.
type DomainResponse struct {
	PublicCertificate interface{} `json:"publicCertificate,omitempty"`
	ValidationStatus interface{} `json:"validationStatus,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// The ID number of the brand
	BrandId string `json:"brandId,omitempty"`
	CertificateSourceType interface{} `json:"certificateSourceType,omitempty"`
	DnsRecords []interface{} `json:"dnsRecords,omitempty"`
	// Custom domain name
	Domain string `json:"domain,omitempty"`
	// Unique ID of the domain
	ID string `json:"id,omitempty"`
}
