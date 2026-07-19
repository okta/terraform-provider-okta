// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EmailDomainResponse represents the EmailDomainResponse schema
type EmailDomainResponse struct {
	DnsValidationRecords []interface{} `json:"dnsValidationRecords,omitempty"`
	Domain string `json:"domain,omitempty"`
	ID string `json:"id,omitempty"`
	ValidationStatus interface{} `json:"validationStatus,omitempty"`
	// The subdomain for the email sender's custom mail domain
	ValidationSubdomain string `json:"validationSubdomain,omitempty"`
}
