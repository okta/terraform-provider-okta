// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OINSaml20ApplicationSettingsSignOn represents the OINSaml20ApplicationSettingsSignOn schema
// Contains SAML 2.0 sign-on mode attributes. > **Note:** Set `destinationOverride` to configure any other SAML 2.0 attributes in this section.
type OINSaml20ApplicationSettingsSignOn struct {
	// A list of custom attribute statements for the app's SAML assertion. See [SAML 2.0 Technical Overview](https://docs.oasis-open.org/security/saml/Post2.0/sstc-saml-tech-overview-2.0-cd-02.html).  The...
	AttributeStatements []interface{} `json:"attributeStatements,omitempty"`
	// Audience override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	AudienceOverride string `json:"audienceOverride,omitempty"`
	// Identifies a specific application resource in an IdP-initiated SSO scenario
	DefaultRelayState string `json:"defaultRelayState,omitempty"`
	// Destination override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	DestinationOverride string `json:"destinationOverride,omitempty"`
	// Recipient override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	RecipientOverride string `json:"recipientOverride,omitempty"`
	// Determines the SAML app session lifetimes with Okta
	SamlAssertionLifetimeSeconds int `json:"samlAssertionLifetimeSeconds,omitempty"`
	// Assertion Consumer Service (ACS) URL override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	SsoAcsUrlOverride string `json:"ssoAcsUrlOverride,omitempty"`
}
