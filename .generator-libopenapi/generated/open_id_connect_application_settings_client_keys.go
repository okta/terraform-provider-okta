// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OpenIdConnectApplicationSettingsClientKeys represents the OpenIdConnectApplicationSettingsClientKeys schema
// A [JSON Web Key Set](https://tools.ietf.org/html/rfc7517#section-5) for validating JWTs presented to Okta or for encrypting ID tokens minted by Okta for the client
type OpenIdConnectApplicationSettingsClientKeys struct {
	Keys []interface{} `json:"keys,omitempty"`
}
