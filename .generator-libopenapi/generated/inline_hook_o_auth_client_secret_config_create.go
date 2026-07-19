// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookOAuthClientSecretConfigCreate represents the InlineHookOAuthClientSecretConfigCreate schema
type InlineHookOAuthClientSecretConfigCreate struct {
	// A private value provided by the service used to authenticate the identity of the app to the service
	ClientSecret string `json:"clientSecret,omitempty"`
	// The method of the Okta inline hook request. Only accepts `POST`.
	Method string `json:"method,omitempty"`
}
