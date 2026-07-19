// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventsProviderSettingsNonSSFCompliant represents the SecurityEventsProviderSettingsNonSSFCompliant schema
// Security events provider with issuer and JWKS settings for signal ingestion
type SecurityEventsProviderSettingsNonSSFCompliant struct {
	// Issuer URL
	Issuer string `json:"issuer"`
	// The public URL where the JWKS public key is uploaded
	JwksUrl string `json:"jwks_url"`
}
