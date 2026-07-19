// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlAcsEndpoint represents the SamlAcsEndpoint schema
// Okta's `SPSSODescriptor` endpoint where the IdP sends a `<SAMLResponse>` message
type SamlAcsEndpoint struct {
	Type interface{} `json:"type,omitempty"`
	Binding interface{} `json:"binding,omitempty"`
}
