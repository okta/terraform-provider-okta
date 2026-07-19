// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlEndpoints represents the SamlEndpoints schema
// SAML 2.0 HTTP binding settings for IdP and SP (Okta)
type SamlEndpoints struct {
	Acs interface{} `json:"acs,omitempty"`
	Slo interface{} `json:"slo,omitempty"`
	Sso interface{} `json:"sso,omitempty"`
}
