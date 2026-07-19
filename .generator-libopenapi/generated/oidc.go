// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Oidc represents the Oidc schema
// OIDC configuration details
type Oidc struct {
	// List of sign-in redirect URIs
	RedirectUris []string `json:"redirectUris"`
	// The URL to your customer-facing instructions for configuring your OIDC integration. See [Customer configuration document guidelines](https://developer.okta.com/docs/guides/submit-app-prereq/main/#c...
	Doc string `json:"doc"`
	// The URL to redirect users when they click on your app from their Okta End-User Dashboard
	InitiateLoginUri string `json:"initiateLoginUri,omitempty"`
	// The sign-out redirect URIs for your app. You can send a request to `/v1/logout` to sign the user out and redirect them to one of these URIs.
	PostLogoutUris []string `json:"postLogoutUris,omitempty"`
}
