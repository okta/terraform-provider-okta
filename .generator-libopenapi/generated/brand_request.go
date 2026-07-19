// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BrandRequest represents the BrandRequest schema
type BrandRequest struct {
	// Consent for updating the custom privacy URL. Not required when resetting the URL.
	AgreeToCustomPrivacyPolicy bool `json:"agreeToCustomPrivacyPolicy,omitempty"`
	// Custom privacy policy URL
	CustomPrivacyPolicyUrl string `json:"customPrivacyPolicyUrl,omitempty"`
	DefaultApp interface{} `json:"defaultApp,omitempty"`
	// The ID of the email domain
	EmailDomainId string `json:"emailDomainId,omitempty"`
	Locale interface{} `json:"locale,omitempty"`
	// The name of the brand  > **Note:** You can't use the reserved `DRAPP_DOMAIN_BRAND` name.
	Name string `json:"name"`
	// Removes "Powered by Okta" from the sign-in page in redirect authentication deployments, and "© [current year] Okta, Inc." from the Okta End-User Dashboard
	RemovePoweredByOkta bool `json:"removePoweredByOkta,omitempty"`
}
