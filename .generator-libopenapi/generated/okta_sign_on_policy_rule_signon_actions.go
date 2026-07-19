// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaSignOnPolicyRuleSignonActions represents the OktaSignOnPolicyRuleSignonActions schema
// Specifies settings for the policy rule
type OktaSignOnPolicyRuleSignonActions struct {
	// Indicates if a user is allowed to sign in
	Access string `json:"access,omitempty"`
	// Interval of time that must elapse before the user is challenged for MFA, if the factor prompt mode is set to `SESSION`  > **Note:** Required only if `requireFactor` is `true`.
	FactorLifetime int `json:"factorLifetime,omitempty"`
	FactorPromptMode interface{} `json:"factorPromptMode,omitempty"`
	PrimaryFactor interface{} `json:"primaryFactor,omitempty"`
	// Indicates if Okta should automatically remember the device
	RememberDeviceByDefault bool `json:"rememberDeviceByDefault,omitempty"`
	// Indicates if multifactor authentication is required
	RequireFactor bool `json:"requireFactor,omitempty"`
	Session interface{} `json:"session,omitempty"`
}
