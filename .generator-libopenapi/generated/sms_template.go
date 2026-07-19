// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// SmsTemplate represents the SmsTemplate schema
type SmsTemplate struct {
	// Text of the Template, including any [macros](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/Template/)
	Template string `json:"template,omitempty"`
	Translations interface{} `json:"translations,omitempty"`
	Type interface{} `json:"type,omitempty"`
	Created *time.Time `json:"created,omitempty"`
	ID string `json:"id,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Human-readable name of the Template
	Name string `json:"name,omitempty"`
}
