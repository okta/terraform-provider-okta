// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OINSaml11ApplicationSettingsSignOn represents the OINSaml11ApplicationSettingsSignOn schema
// Contains SAML 1.1 sign-on mode attributes
type OINSaml11ApplicationSettingsSignOn struct {
	// Assertion Consumer Service (ACS) URL override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	SsoAcsUrlOverride string `json:"ssoAcsUrlOverride,omitempty"`
	// Audience override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	AudienceOverride string `json:"audienceOverride,omitempty"`
	// Identifies a specific application resource in an IdP-initiated SSO scenario
	DefaultRelayState string `json:"defaultRelayState,omitempty"`
	// Recipient override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	RecipientOverride string `json:"recipientOverride,omitempty"`
}
