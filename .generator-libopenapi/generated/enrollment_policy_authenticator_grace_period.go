// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EnrollmentPolicyAuthenticatorGracePeriod represents the EnrollmentPolicyAuthenticatorGracePeriod schema
// Specifies the [grace period](https://developer.okta.com/docs/concepts/policies/#authenticator-enrollment-policies) configuration for completing an authenticator enrollment or setup
type EnrollmentPolicyAuthenticatorGracePeriod struct {
	// Grace period type  * `BY_DATE_TIME`: The grace period is defined by a specific date and time. * <x-lifecycle class="ea"></x-lifecycle>`BY_SKIP_COUNT`: The grace period is defined by the number of t...
	Type string `json:"type,omitempty"`
}
