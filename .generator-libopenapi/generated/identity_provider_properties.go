// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdentityProviderProperties represents the IdentityProviderProperties schema
// The properties in the IdP `properties` object vary depending on the IdP type
type IdentityProviderProperties struct {
	// The [authentication assurance level](https://developers.login.gov/oidc/#aal-values) (AAL) value for the Login.gov IdP. See [Add a Login.gov IdP](https://developer.okta.com/docs/guides/add-logingov-...
	AalValue string `json:"aalValue,omitempty"`
	// The additional Assurance Methods References (AMR) values for Smart Card IdPs. Applies to `X509` IdP type.
	AdditionalAmr []string `json:"additionalAmr,omitempty"`
	// The [type of identity verification](https://developers.login.gov/oidc/#ial-values) (IAL) value for the Login.gov IdP. See [Add a Login.gov IdP](https://developer.okta.com/docs/guides/add-logingov-i...
	IalValue string `json:"ialValue,omitempty"`
	// Metadata about the IDV vendor. Available only for `IDV_STANDARD` IdPs.
	IdvMetadata map[string]interface{} `json:"idvMetadata,omitempty"`
	// The ID of the inquiry template from your Persona dashboard. The inquiry template always starts with `itmpl`. Applies to the `IDV_PERSONA` IdP type.
	InquiryTemplateId string `json:"inquiryTemplateId"`
}
