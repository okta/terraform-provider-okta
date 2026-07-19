// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// WellKnownSSFMetadata represents the WellKnownSSFMetadata schema
// Metadata about Okta as a transmitter and relevant information for configuration.
type WellKnownSSFMetadata struct {
	// An array of supported SET delivery methods
	DeliveryMethodsSupported []string `json:"delivery_methods_supported,omitempty"`
	// The issuer used in security event tokens. This value is set as `iss` in the claim.
	Issuer string `json:"issuer,omitempty"`
	// The URL of the JSON Web Key Set (JWKS) that contains the signing keys for validating the signatures of security event tokens (SETs)
	JwksUri string `json:"jwks_uri,omitempty"`
	// The version identifying the implementer's draft or final specification implemented by the transmitter
	SpecVersion string `json:"spec_version,omitempty"`
	// The URL of the SSF stream verification endpoint
	VerificationEndpoint string `json:"verification_endpoint,omitempty"`
	// An array of JSON objects that specify the authorization scheme properties supported by the transmitter
	AuthorizationSchemes []interface{} `json:"authorization_schemes,omitempty"`
	// The URL of the SSF stream configuration endpoint
	ConfigurationEndpoint string `json:"configuration_endpoint,omitempty"`
	// A string that indicates the default behavior of newly created streams
	DefaultSubjects string `json:"default_subjects,omitempty"`
}
