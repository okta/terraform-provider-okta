// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdentitySourceUserProfileForUpsert represents the IdentitySourceUserProfileForUpsert schema
// Contains a set of external user attributes and their values that are mapped to Okta standard and custom profile properties. See the [`profile` object](https://developer.okta.com/docs/api/openapi/ok...
type IdentitySourceUserProfileForUpsert struct {
	// Alternative email address of the user
	SecondEmail string `json:"secondEmail,omitempty"`
	// Username of the user
	UserName string `json:"userName,omitempty"`
	// Email address of the user
	Email string `json:"email,omitempty"`
	// First name of the user
	FirstName string `json:"firstName,omitempty"`
	// Home address of the user
	HomeAddress string `json:"homeAddress,omitempty"`
	// Last name of the user
	LastName string `json:"lastName,omitempty"`
	// Mobile phone number of the user
	MobilePhone string `json:"mobilePhone,omitempty"`
}
