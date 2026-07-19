// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UserActivationToken represents the UserActivationToken schema
type UserActivationToken struct {
	// Token received as part of an activation user request. If a password was set before the user was activated, then user must sign in with their password or the `activationToken` and not the activation...
	ActivationToken string `json:"activationToken,omitempty"`
	// If `sendEmail` is `false`, returns an activation link for the user to set up their account. You can use the activation token to create a custom activation link.  > **Note:** The `activationUrl` var...
	ActivationUrl string `json:"activationUrl,omitempty"`
}
