// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// OrgOktaSupportSettingsObj represents the OrgOktaSupportSettingsObj schema
type OrgOktaSupportSettingsObj struct {
	Support interface{} `json:"support,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for the Okta Support Settings object using the [JSON Hypertext Application Language](https://datatrack...
	Links map[string]interface{} `json:"_links,omitempty"`
	// Support case number for the Okta Support access grant
	CaseNumber string `json:"caseNumber,omitempty"`
	// Expiration of Okta Support
	Expiration *time.Time `json:"expiration,omitempty"`
}
