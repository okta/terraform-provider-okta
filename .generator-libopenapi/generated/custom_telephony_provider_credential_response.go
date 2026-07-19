// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderCredentialResponse represents the CustomTelephonyProviderCredentialResponse schema
type CustomTelephonyProviderCredentialResponse struct {
	// Indicates whether the provider is enabled and can be used
	Enabled bool `json:"enabled,omitempty"`
	// ID of the custom telephony provider
	ID string `json:"id,omitempty"`
	// Indicates whether the provider is the primary telephony provider
	IsPrimaryProvider bool `json:"isPrimaryProvider,omitempty"`
	// The types of telephony operations (SMS or Voice) that you use with your telephony provider.  `ALL` is the only valid value. It indicates that your provider can handle both SMS messages and voice ca...
	ProviderCapability string `json:"providerCapability,omitempty"`
	// Name of the telephony provider
	ProviderName string `json:"providerName,omitempty"`
	ProviderSettings interface{} `json:"providerSettings,omitempty"`
	// The account string identifier (SID) for your telephony provider account. Your telephony provider gives you this SID.
	ProviderSid string `json:"providerSid,omitempty"`
}
