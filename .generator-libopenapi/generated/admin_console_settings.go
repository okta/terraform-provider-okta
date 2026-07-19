// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AdminConsoleSettings represents the AdminConsoleSettings schema
// Settings specific to the Okta Admin Console
type AdminConsoleSettings struct {
	// The absolute maximum session lifetime of the Okta Admin Console. Must be no more than 7 days.
	SessionMaxLifetimeMinutes int `json:"sessionMaxLifetimeMinutes,omitempty"`
	// The maximum idle time before the Okta Admin Console session expires. Must be no more than 12 hours.
	SessionIdleTimeoutMinutes int `json:"sessionIdleTimeoutMinutes,omitempty"`
}
