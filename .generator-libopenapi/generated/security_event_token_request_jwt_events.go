// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SecurityEventTokenRequestJwtEvents represents the SecurityEventTokenRequestJwtEvents schema
// A non-empty collection of events
type SecurityEventTokenRequestJwtEvents struct {
	Https://schemas.okta.com/secevent/okta/event-type/device-risk-change interface{} `json:"https://schemas.okta.com/secevent/okta/event-type/device-risk-change,omitempty"`
	Https://schemas.okta.com/secevent/okta/event-type/ip-change interface{} `json:"https://schemas.okta.com/secevent/okta/event-type/ip-change,omitempty"`
	Https://schemas.okta.com/secevent/okta/event-type/user-risk-change interface{} `json:"https://schemas.okta.com/secevent/okta/event-type/user-risk-change,omitempty"`
	Https://schemas.openid.net/secevent/caep/event-type/device-compliance-change interface{} `json:"https://schemas.openid.net/secevent/caep/event-type/device-compliance-change,omitempty"`
	Https://schemas.openid.net/secevent/caep/event-type/session-revoked interface{} `json:"https://schemas.openid.net/secevent/caep/event-type/session-revoked,omitempty"`
	Https://schemas.openid.net/secevent/risc/event-type/identifier-changed interface{} `json:"https://schemas.openid.net/secevent/risc/event-type/identifier-changed,omitempty"`
}
