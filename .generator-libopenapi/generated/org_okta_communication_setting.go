// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OrgOktaCommunicationSetting represents the OrgOktaCommunicationSetting schema
type OrgOktaCommunicationSetting struct {
	// Indicates whether org users receive Okta communication emails
	OptOutEmailUsers bool `json:"optOutEmailUsers,omitempty"`
	// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available for this object using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/...
	Links map[string]interface{} `json:"_links,omitempty"`
}
