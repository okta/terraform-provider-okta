// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventsProviderSettingsResponse represents the SecurityEventsProviderSettingsResponse schema
// Security events provider settings
type SecurityEventsProviderSettingsResponse struct {
	// Issuer URL
	Issuer string `json:"issuer,omitempty"`
	// The public URL where the JWKS public key is uploaded
	JwksUrl string `json:"jwks_url,omitempty"`
	// The well-known URL of the security events provider (the SSF transmitter)
	WellKnownUrl string `json:"well_known_url,omitempty"`
}
