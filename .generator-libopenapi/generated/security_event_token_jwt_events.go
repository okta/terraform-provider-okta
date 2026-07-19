// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventTokenJwtEvents represents the SecurityEventTokenJwtEvents schema
// A non-empty set of events. Expected size is 1 for each SET
type SecurityEventTokenJwtEvents struct {
	Https://schemas.openid.net/secevent/caep/event-type/credential-change interface{} `json:"https://schemas.openid.net/secevent/caep/event-type/credential-change,omitempty"`
	Https://schemas.openid.net/secevent/caep/event-type/session-revoked interface{} `json:"https://schemas.openid.net/secevent/caep/event-type/session-revoked,omitempty"`
}
