// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticatorEnrollmentPolicySettings represents the AuthenticatorEnrollmentPolicySettings schema
// Specifies the policy level settings  > **Note:** In Identity Engine, the Multifactor (MFA) Enrollment policy name has changed to authenticator enrollment policy. The policy type of `MFA_ENROLL` rem...
type AuthenticatorEnrollmentPolicySettings struct {
	// List of authenticator policy settings  <x-lifecycle class="oie"></x-lifecycle> For orgs with the Authenticator enrollment policy feature enabled, the new default authenticator enrollment policy cre...
	Authenticators []interface{} `json:"authenticators,omitempty"`
	Type interface{} `json:"type,omitempty"`
}
