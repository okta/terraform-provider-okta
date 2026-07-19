// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EmailPreview represents the EmailPreview schema
type EmailPreview struct {
	// The email's HTML body
	Body string `json:"body,omitempty"`
	// The email's subject
	Subject string `json:"subject,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
