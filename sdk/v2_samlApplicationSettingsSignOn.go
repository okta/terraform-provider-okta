package sdk

type SamlApplicationSettingsSignOn struct {
	AcsEndpoints              []*AcsEndpoint            `json:"acsEndpoints,omitempty"`
	AllowMultipleAcsEndpoints *bool                     `json:"allowMultipleAcsEndpoints,omitempty"`
	AssertionSigned           *bool                     `json:"assertionSigned,omitempty"`
	AttributeStatements       []*SamlAttributeStatement `json:"attributeStatements"`
	Audience                  string                    `json:"audience,omitempty"`
	AudienceOverride          string                    `json:"audienceOverride"`
	AuthnContextClassRef      string                    `json:"authnContextClassRef,omitempty"`
	DefaultRelayState         string                    `json:"defaultRelayState"`
	Destination               string                    `json:"destination,omitempty"`
	DestinationOverride       string                    `json:"destinationOverride"`
	DigestAlgorithm           string                    `json:"digestAlgorithm,omitempty"`
	HonorForceAuthn           *bool                     `json:"honorForceAuthn,omitempty"`
	IdpIssuer                 string                    `json:"idpIssuer,omitempty"`
	InlineHooks               []*SignOnInlineHook       `json:"inlineHooks,omitempty"`
	Recipient                 string                    `json:"recipient,omitempty"`
	RecipientOverride         string                    `json:"recipientOverride"`
	RequestCompressed         *bool                     `json:"requestCompressed,omitempty"`
	ResponseSigned            *bool                     `json:"responseSigned,omitempty"`
	SamlSignedRequestEnabled  *bool                     `json:"samlSignedRequestEnabled,omitempty"`
	SignatureAlgorithm        string                    `json:"signatureAlgorithm,omitempty"`
	Slo                       *SingleLogout             `json:"slo,omitempty"`
	SpCertificate             *SpCertificate            `json:"spCertificate,omitempty"`
	SpIssuer                  string                    `json:"spIssuer,omitempty"`
	SsoAcsUrl                 string                    `json:"ssoAcsUrl,omitempty"`
	SsoAcsUrlOverride         string                    `json:"ssoAcsUrlOverride"`
	SubjectNameIdFormat       string                    `json:"subjectNameIdFormat,omitempty"`
	SubjectNameIdTemplate     string                    `json:"subjectNameIdTemplate,omitempty"`
}
