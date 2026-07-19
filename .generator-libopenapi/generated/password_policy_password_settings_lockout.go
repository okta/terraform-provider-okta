// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyPasswordSettingsLockout represents the PasswordPolicyPasswordSettingsLockout schema
// Lockout settings
type PasswordPolicyPasswordSettingsLockout struct {
	// Specifies the time interval (in minutes) a locked account remains locked before it is automatically unlocked: `0` indicates no limit
	AutoUnlockMinutes int `json:"autoUnlockMinutes,omitempty"`
	// Specifies the number of times Users can attempt to sign in to their accounts with an invalid password before their accounts are locked: `0` indicates no limit
	MaxAttempts int `json:"maxAttempts,omitempty"`
	// Indicates if the User should be informed when their account is locked
	ShowLockoutFailures bool `json:"showLockoutFailures,omitempty"`
	// How the user is notified when their account becomes locked. The only acceptable values are `[]` and `['EMAIL']`.
	UserLockoutNotificationChannels []string `json:"userLockoutNotificationChannels,omitempty"`
}
