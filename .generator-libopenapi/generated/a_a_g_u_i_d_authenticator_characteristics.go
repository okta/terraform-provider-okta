// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AAGUIDAuthenticatorCharacteristics represents the AAGUIDAuthenticatorCharacteristics schema
// Contains additional properties about custom AAGUID.
type AAGUIDAuthenticatorCharacteristics struct {
	// Indicates whether the authenticator meets Federal Information Processing Standards (FIPS) compliance requirements
	FipsCompliant bool `json:"fipsCompliant,omitempty"`
	// Indicates whether the authenticator stores the private key on a hardware component
	HardwareProtected bool `json:"hardwareProtected,omitempty"`
	// Indicates whether the custom AAGUID is built into the authenticator (`true`) or if it's a separate, external authenticator
	PlatformAttached bool `json:"platformAttached,omitempty"`
}
