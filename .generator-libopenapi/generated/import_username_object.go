// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ImportUsernameObject represents the ImportUsernameObject schema
// Determines the Okta username for the imported user
type ImportUsernameObject struct {
	// For `usernameFormat=CUSTOM`, specifies the Okta Expression Language statement for a username format that imported users use to sign in to Okta
	UserNameExpression string `json:"userNameExpression,omitempty"`
	// Determines the username format when users sign in to Okta
	UsernameFormat string `json:"usernameFormat"`
}
