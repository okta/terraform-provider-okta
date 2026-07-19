// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EmailDomain represents the EmailDomain schema
type EmailDomain struct {
	// Subdomain for the email sender's custom mail domain. Specify your subdomain when you configure a custom mail domain.
	ValidationSubdomain string `json:"validationSubdomain,omitempty"`
	BrandId string `json:"brandId"`
	Domain string `json:"domain"`
}
