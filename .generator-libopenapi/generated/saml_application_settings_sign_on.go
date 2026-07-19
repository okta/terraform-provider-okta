// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlApplicationSettingsSignOn represents the SamlApplicationSettingsSignOn schema
// SAML 2.0 sign-on attributes. > **Note:** Set either `destinationOverride` or `ssoAcsUrl` to configure any other SAML 2.0 attributes in this section.
type SamlApplicationSettingsSignOn struct {
	// The list of dynamic attribute statements for the SAML assertion inherited from app metadata (apps from the OIN) during app creation.  There are two types of attribute statements: `EXPRESSION` and `...
	ConfiguredAttributeStatements []interface{} `json:"configuredAttributeStatements,omitempty"`
	// Determines the digest algorithm used to digitally sign the SAML assertion and response
	DigestAlgorithm string `json:"digestAlgorithm"`
	// Determines whether the SAML authentication response message is digitally signed by the IdP > **Note:** Either (or both) `responseSigned` or `assertionSigned` must be `TRUE`.
	ResponseSigned bool `json:"responseSigned"`
	// Determines the signing algorithm used to digitally sign the SAML assertion and response
	SignatureAlgorithm string `json:"signatureAlgorithm"`
	// Assertion Consumer Service (ACS) URL override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	SsoAcsUrlOverride string `json:"ssoAcsUrlOverride,omitempty"`
	// Determines whether the app allows you to configure multiple ACS URIs
	AllowMultipleAcsEndpoints bool `json:"allowMultipleAcsEndpoints"`
	// A list of custom attribute statements for the app's SAML assertion. See [SAML 2.0 Technical Overview](https://docs.oasis-open.org/security/saml/Post2.0/sstc-saml-tech-overview-2.0-cd-02.html).  The...
	AttributeStatements []interface{} `json:"attributeStatements,omitempty"`
	// Identifies the location inside the SAML assertion where the SAML response should be sent
	Destination string `json:"destination"`
	// SAML Issuer ID
	IdpIssuer string `json:"idpIssuer"`
	// Identifies the SAML authentication context class for the assertion's authentication statement
	AuthnContextClassRef string `json:"authnContextClassRef"`
	// Identifies a specific application resource in an IdP-initiated SSO scenario
	DefaultRelayState string `json:"defaultRelayState,omitempty"`
	// Destination override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	DestinationOverride string `json:"destinationOverride,omitempty"`
	Slo interface{} `json:"slo,omitempty"`
	// Associates the app with SAML inline hooks. See [the SAML assertion inline hook reference](https://developer.okta.com/docs/reference/saml-hook/).
	InlineHooks []interface{} `json:"inlineHooks,omitempty"`
	// The location where the app may present the SAML assertion
	Recipient string `json:"recipient"`
	// Determines the SAML app session lifetimes with Okta
	SamlAssertionLifetimeSeconds int `json:"samlAssertionLifetimeSeconds,omitempty"`
	// The issuer ID for the Service Provider. This property appears when SLO is enabled.
	SpIssuer string `json:"spIssuer,omitempty"`
	// Identifies the SAML processing rules. Supported values:
	SubjectNameIdFormat string `json:"subjectNameIdFormat"`
	AssertionEncryption interface{} `json:"assertionEncryption,omitempty"`
	ParticipateSlo interface{} `json:"participateSlo,omitempty"`
	// Recipient override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	RecipientOverride string `json:"recipientOverride,omitempty"`
	// Determines whether the SAML request is expected to be compressed
	RequestCompressed bool `json:"requestCompressed"`
	// Template for app user's username when a user is assigned to the app
	SubjectNameIdTemplate string `json:"subjectNameIdTemplate"`
	// Determines whether the SAML assertion is digitally signed
	AssertionSigned bool `json:"assertionSigned"`
	// Set to `true` to prompt users for their credentials when a SAML request has the `ForceAuthn` attribute set to `true`
	HonorForceAuthn bool `json:"honorForceAuthn"`
	SpCertificate interface{} `json:"spCertificate,omitempty"`
	// An array of ACS endpoints. You can configure a maximum of 100 endpoints.
	AcsEndpoints []interface{} `json:"acsEndpoints,omitempty"`
	// Single Sign-On Assertion Consumer Service (ACS) URL
	SsoAcsUrl string `json:"ssoAcsUrl"`
	// The entity ID of the SP. Use the entity ID value exactly as provided by the SP.
	Audience string `json:"audience"`
	// Audience override for CASB configuration. See [CASB config guide](https://help.okta.com/en-us/Content/Topics/Apps/CASB-config-guide.htm).
	AudienceOverride string `json:"audienceOverride,omitempty"`
}
