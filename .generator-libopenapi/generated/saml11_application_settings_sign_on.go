// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Saml11ApplicationSettingsSignOn represents the Saml11ApplicationSettingsSignOn schema
// SAML 1.1 sign-on mode attributes
type Saml11ApplicationSettingsSignOn struct {
	// The intended audience of the SAML assertion. This is usually the Entity ID of your application.
	AudienceOverride string `json:"audienceOverride,omitempty"`
	// The URL of the resource to direct users after they successfully sign in to the SP using SAML. See the SP documentation to check if you need to specify a RelayState. In most instances, you can leave...
	DefaultRelayState string `json:"defaultRelayState,omitempty"`
	// The location where the application can present the SAML assertion. This is usually the Single Sign-On (SSO) URL.
	RecipientOverride string `json:"recipientOverride,omitempty"`
	// Assertion Consumer Services (ACS) URL value for the Service Provider (SP). This URL is always used for Identity Provider (IdP) initiated sign-on requests.
	SsoAcsUrlOverride string `json:"ssoAcsUrlOverride,omitempty"`
}
