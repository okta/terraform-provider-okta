// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppUserCredentials represents the AppUserCredentials schema
// Specifies a user's credentials for the app. This parameter can be omitted for apps with [sign-on mode](/openapi/okta-management/management/application/getapplication#application/getapplication/t=re...
type AppUserCredentials struct {
	Password interface{} `json:"password,omitempty"`
	// The user's username in the app  > **Note:** The [userNameTemplate](/openapi/okta-management/management/tags/application/other/createapplication#application/createapplication/t=response&c=200&path=&...
	UserName string `json:"userName,omitempty"`
}
