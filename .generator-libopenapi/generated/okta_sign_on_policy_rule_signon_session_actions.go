// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaSignOnPolicyRuleSignonSessionActions represents the OktaSignOnPolicyRuleSignonSessionActions schema
// Properties governing the user's session lifetime
type OktaSignOnPolicyRuleSignonSessionActions struct {
	// Maximum number of minutes that a user session can be idle before the session is ended
	MaxSessionIdleMinutes int `json:"maxSessionIdleMinutes,omitempty"`
	// Maximum number of minutes (from when the user signs in) that a user's session is active. Set this to force users to sign in again after the number of specified minutes. Disable by setting to `0`.
	MaxSessionLifetimeMinutes int `json:"maxSessionLifetimeMinutes,omitempty"`
	// If set to `false`, user session cookies only last the length of a browser session. If set to `true`, user session cookies last across browser sessions. This setting doesn't impact administrators wh...
	UsePersistentCookie bool `json:"usePersistentCookie,omitempty"`
}
