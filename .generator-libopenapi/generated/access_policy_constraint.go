// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AccessPolicyConstraint represents the AccessPolicyConstraint schema
// Consists of a `POSSESSION` constraint, a `KNOWLEDGE` constraint, or both. You can't configure an `INHERENCE` constraint, but an inherence factor can satisfy the second part of a 2FA assurance if no...
type AccessPolicyConstraint struct {
	// This property specifies the precise authenticator and method to exclude from authentication. <x-lifecycle class="oie"></x-lifecycle>
	ExcludedAuthenticationMethods interface{} `json:"excludedAuthenticationMethods,omitempty"`
	// The authenticator methods that are permitted
	Methods []string `json:"methods,omitempty"`
	// The duration after which the user must re-authenticate regardless of user activity. This re-authentication interval overrides the Verification Method object's `reauthenticateIn` interval. The suppo...
	ReauthenticateIn string `json:"reauthenticateIn,omitempty"`
	// This property indicates whether the knowledge or possession factor is required by the assurance. It's optional in the request, but is always returned in the response. By default, this field is `tru...
	Required bool `json:"required,omitempty"`
	// The authenticator types that are permitted
	Types []string `json:"types,omitempty"`
	// This property specifies the precise authenticator and method for authentication. <x-lifecycle class="oie"></x-lifecycle>
	AuthenticationMethods []interface{} `json:"authenticationMethods,omitempty"`
}
