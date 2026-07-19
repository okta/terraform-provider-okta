// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Brand represents the Brand schema
type Brand struct {
	// Removes "Powered by Okta" from the sign-in page in redirect authentication deployments, and "© [current year] Okta, Inc." from the Okta End-User Dashboard
	RemovePoweredByOkta bool `json:"removePoweredByOkta,omitempty"`
	// Consent for updating the custom privacy URL. Not required when resetting the URL.
	AgreeToCustomPrivacyPolicy bool `json:"agreeToCustomPrivacyPolicy,omitempty"`
	DefaultApp interface{} `json:"defaultApp,omitempty"`
	Locale interface{} `json:"locale,omitempty"`
	// Custom privacy policy URL
	CustomPrivacyPolicyUrl string `json:"customPrivacyPolicyUrl,omitempty"`
	// The ID of the email domain
	EmailDomainId string `json:"emailDomainId,omitempty"`
	// The Brand ID
	ID string `json:"id,omitempty"`
	// If `true`, the Brand is used for the Okta subdomain
	IsDefault bool `json:"isDefault,omitempty"`
	// The name of the Brand
	Name string `json:"name,omitempty"`
}
