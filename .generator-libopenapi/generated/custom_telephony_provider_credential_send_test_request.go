// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderCredentialSendTestRequest represents the CustomTelephonyProviderCredentialSendTestRequest schema
type CustomTelephonyProviderCredentialSendTestRequest struct {
	// The type of test message to send
	Factor string `json:"factor,omitempty"`
	// The phone number to which the test message or call is sent
	PhoneNumber string `json:"phoneNumber,omitempty"`
	// The country code for the phone number. Use the [Alpha-2 code from ISO 3166-1](https://www.iso.org/obp/ui/#search) for country codes.
	CountryCodeIso2 string `json:"countryCodeIso2,omitempty"`
}
