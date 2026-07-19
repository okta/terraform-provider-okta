// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordPolicyPasswordSettingsAge represents the PasswordPolicyPasswordSettingsAge schema
// Age settings
type PasswordPolicyPasswordSettingsAge struct {
	// Specifies the number of days prior to password expiration when a User is warned to reset their password: `0` indicates no warning
	ExpireWarnDays int `json:"expireWarnDays,omitempty"`
	// Specifies the number of distinct passwords that a User must create before they can reuse a previous password: `0` indicates none
	HistoryCount int `json:"historyCount,omitempty"`
	// Specifies how long (in days) a password remains valid before it expires: `0` indicates no limit
	MaxAgeDays int `json:"maxAgeDays,omitempty"`
	// Specifies the minimum time interval (in minutes) between password changes: `0` indicates no limit
	MinAgeMinutes int `json:"minAgeMinutes,omitempty"`
}
