// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SAMLPayLoad represents the SAMLPayLoad schema
type SAMLPayLoad struct {
	Data map[string]interface{} `json:"data,omitempty"`
	// The type of inline hook. The SAML assertion inline hook type is `com.okta.saml.tokens.transform`.
	EventType string `json:"eventType,omitempty"`
	// The ID and URL of the SAML assertion inline hook
	Source string `json:"source,omitempty"`
}
