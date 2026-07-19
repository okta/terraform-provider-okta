// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OidcSettings represents the OidcSettings schema
// Advanced settings for the OpenID Connect protocol
type OidcSettings struct {
	// Determines if the IdP should send the application context as `OktaAppInstanceId` and `OktaAppName` params in the request
	SendApplicationContext bool `json:"sendApplicationContext,omitempty"`
	// Set to `true` to have Okta send a logout request to the upstream IdP when a user signs out of Okta or a downstream app.
	ParticipateSlo bool `json:"participateSlo,omitempty"`
}
