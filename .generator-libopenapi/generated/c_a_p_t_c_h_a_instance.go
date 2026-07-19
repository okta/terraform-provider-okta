// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CAPTCHAInstance represents the CAPTCHAInstance schema
type CAPTCHAInstance struct {
	// The unique key for the CAPTCHA instance
	ID string `json:"id,omitempty"`
	// The name of the CAPTCHA instance
	Name string `json:"name,omitempty"`
	// The secret key issued from the CAPTCHA provider to perform server-side validation for a CAPTCHA token
	SecretKey string `json:"secretKey,omitempty"`
	// The site key issued from the CAPTCHA provider to render a CAPTCHA on a page
	SiteKey string `json:"siteKey,omitempty"`
	Type interface{} `json:"type,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
