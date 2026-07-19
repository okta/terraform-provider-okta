// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CaepCredentialChangeEvent represents the CaepCredentialChangeEvent schema
// The credential was created, changed, revoked or deleted
type CaepCredentialChangeEvent struct {
	// The credential type of the changed credential. It's one of the supported enum values or any other credential type supported mutually by the transmitter and the receiver.
	CredentialType string `json:"credential_type"`
	// The time of the event (UNIX timestamp)
	EventTimestamp int64 `json:"event_timestamp,omitempty"`
	// FIDO2 Authenticator Attestation GUID
	Fido2Aaguid string `json:"fido2_aaguid,omitempty"`
	// The entity that initiated the event
	InitiatingEntity string `json:"initiating_entity,omitempty"`
	ReasonAdmin interface{} `json:"reason_admin,omitempty"`
	ReasonUser interface{} `json:"reason_user,omitempty"`
	// The type of action done towards the credential
	ChangeType string `json:"change_type"`
	// Credential friendly name
	FriendlyName string `json:"friendly_name,omitempty"`
	Subject interface{} `json:"subject,omitempty"`
}
