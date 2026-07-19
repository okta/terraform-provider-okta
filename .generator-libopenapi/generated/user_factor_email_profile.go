// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserFactorEmailProfile represents the UserFactorEmailProfile schema
type UserFactorEmailProfile struct {
	// Email address of the user. This must be either the primary or secondary email address associated with the Okta user account.  > **Note:** For Identity Engine orgs, you can only enroll the primary e...
	Email string `json:"email,omitempty"`
}
