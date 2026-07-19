// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OrgCAPTCHASettings represents the OrgCAPTCHASettings schema
type OrgCAPTCHASettings struct {
	// The unique key of the associated CAPTCHA instance
	CaptchaId string `json:"captchaId,omitempty"`
	// An array of pages that have CAPTCHA enabled
	EnabledPages []interface{} `json:"enabledPages,omitempty"`
	// Link relations for the CAPTCHA settings object
	Links map[string]interface{} `json:"_links,omitempty"`
}
