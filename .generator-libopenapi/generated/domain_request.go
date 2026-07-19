// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DomainRequest represents the DomainRequest schema
type DomainRequest struct {
	CertificateSourceType interface{} `json:"certificateSourceType"`
	// Custom domain name  > **Note:** You can't use the reserved `drapp.{yourOrgSubDomain}.okta.com` domain.
	Domain string `json:"domain"`
}
