// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProfileEnrollmentPolicyRuleAction represents the ProfileEnrollmentPolicyRuleAction schema
type ProfileEnrollmentPolicyRuleAction struct {
	// Indicates if the user profile is granted access  > **Note:** You can't set the `access` property to `DENY` after you create the policy
	Access string `json:"access,omitempty"`
	ActivationRequirements interface{} `json:"activationRequirements,omitempty"`
	// A list of attributes to identify an end user. Can be used across Okta sign-in, unlock, and recovery flows.
	AllowedIdentifiers []string `json:"allowedIdentifiers,omitempty"`
	// Additional authenticator fields that can be used on the first page of user registration. Valid values only includes `'password'`.
	EnrollAuthenticatorTypes []string `json:"enrollAuthenticatorTypes,omitempty"`
	// (Optional) The `id` of at most one registration inline hook
	PreRegistrationInlineHooks []interface{} `json:"preRegistrationInlineHooks,omitempty"`
	// Progressive profile enrollment helps evaluate the user profile policy at every user login. Users can be prompted to provide input for newly required attributes.
	ProgressiveProfilingAction string `json:"progressiveProfilingAction,omitempty"`
	// (Optional, max 1 entry) The `id` of a group that this user should be added to
	TargetGroupIds []string `json:"targetGroupIds,omitempty"`
	// Value created by the backend. If present, all policy updates must include this attribute/value.
	UiSchemaId string `json:"uiSchemaId,omitempty"`
	// A list of attributes to prompt the user for during registration or progressive profiling. Where defined on the user schema, these attributes are persisted in the user profile. You can also add non-...
	ProfileAttributes []interface{} `json:"profileAttributes,omitempty"`
	// Which action should be taken if this user is new
	UnknownUserAction string `json:"unknownUserAction,omitempty"`
}
