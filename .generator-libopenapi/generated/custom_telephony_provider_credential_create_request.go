// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderCredentialCreateRequest represents the CustomTelephonyProviderCredentialCreateRequest schema
// Create custom telephony provider credentials
type CustomTelephonyProviderCredentialCreateRequest struct {
	// The types of telephony operations (SMS or Voice) that you use with your telephony provider.  `ALL` is the only valid value. It indicates that your provider can handle both SMS messages and voice ca...
	ProviderCapability string `json:"providerCapability,omitempty"`
	// The name of the telephony provider
	ProviderName string `json:"providerName,omitempty"`
	ProviderSettings interface{} `json:"providerSettings,omitempty"`
	// The account string identifier (SID) for your telephony provider account. Your telephony provider gives you this SID.
	ProviderSid string `json:"providerSid,omitempty"`
	// The authentication token that's used to authenticate requests to the telephony provider. Your telephony provider gives you this token.
	ProviderAuthToken string `json:"providerAuthToken,omitempty"`
}
