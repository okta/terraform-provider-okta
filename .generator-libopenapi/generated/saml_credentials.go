// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlCredentials represents the SamlCredentials schema
// Federation Trust Credentials for verifying assertions from the IdP and signing requests to the IdP
type SamlCredentials struct {
	Signing interface{} `json:"signing,omitempty"`
	Trust interface{} `json:"trust,omitempty"`
}
