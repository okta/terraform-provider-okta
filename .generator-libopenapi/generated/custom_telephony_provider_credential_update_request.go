// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderCredentialUpdateRequest represents the CustomTelephonyProviderCredentialUpdateRequest schema
// Update custom telephony provider credentials
type CustomTelephonyProviderCredentialUpdateRequest struct {
	// The authentication token that's used to authenticate requests to the telephony provider. Your telephony provider gives you this token.
	ProviderAuthToken string `json:"providerAuthToken,omitempty"`
	ProviderSettings interface{} `json:"providerSettings,omitempty"`
	// The account string identifier (SID) for your telephony provider account. Your telephony provider gives you this SID.
	ProviderSid string `json:"providerSid,omitempty"`
	// ID of the custom telephony provider
	ID string `json:"id,omitempty"`
}
