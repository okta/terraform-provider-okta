// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyPasswordSettingsBreachedProtection represents the PasswordPolicyPasswordSettingsBreachedProtection schema
// Breached Protection settings
type PasswordPolicyPasswordSettingsBreachedProtection struct {
	// The `id` of the workflow that runs when a breached password is found during a sign-in attempt.
	DelegatedWorkflowId string `json:"delegatedWorkflowId,omitempty"`
	// Specifies the number of days after a breached password is found during a sign-in attempt that the user's password should expire. Valid values: 0 through 10. If set to 0, it happens immediately.
	ExpireAfterDays int `json:"expireAfterDays,omitempty"`
	// (Optional, default is false) If true, you must also specify a value for `expireAfterDays`. When enabled, the user's session(s) are terminated immediately the first time the user's credentials are d...
	LogoutEnabled bool `json:"logoutEnabled,omitempty"`
}
