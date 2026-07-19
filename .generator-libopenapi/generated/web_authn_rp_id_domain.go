// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// WebAuthnRpIdDomain represents the WebAuthnRpIdDomain schema
// The RP domain object for the WebAuthn configuration
type WebAuthnRpIdDomain struct {
	// The [RP ID](https://www.w3.org/TR/webauthn/#relying-party-identifier) domain value to be used for all WebAuthn operations.  If it isn't specified, the `domain` object isn't included in the request,...
	Name string `json:"name,omitempty"`
	// Indicates the validation status of the domain
	ValidationStatus string `json:"validationStatus,omitempty"`
	DnsRecord interface{} `json:"dnsRecord,omitempty"`
}
