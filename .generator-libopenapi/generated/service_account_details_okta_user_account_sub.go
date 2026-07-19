// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ServiceAccountDetailsOktaUserAccountSub represents the ServiceAccountDetailsOktaUserAccountSub schema
// Details for managing an Okta user as a service account
type ServiceAccountDetailsOktaUserAccountSub struct {
	Credentials interface{} `json:"credentials,omitempty"`
	// The email address for the Okta user
	Email string `json:"email,omitempty"`
	// The ID of the Okta user to manage as a service account
	OktaUserId string `json:"oktaUserId"`
}
