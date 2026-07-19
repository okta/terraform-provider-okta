// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DNSRecordAuthenticators represents the DNSRecordAuthenticators schema
// DNS TXT record that must be registered for an RP ID domain that requires verification. This is used to verify the domain ownership for the WebAuthn RP ID configuration. After the domain ownership i...
type DNSRecordAuthenticators struct {
	// The DNS record name
	Fqdn string `json:"fqdn,omitempty"`
	RecordType interface{} `json:"recordType,omitempty"`
	// The DNS record verification value
	VerificationValue string `json:"verificationValue,omitempty"`
}
