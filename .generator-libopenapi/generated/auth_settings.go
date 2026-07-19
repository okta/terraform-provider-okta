// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthSettings represents the AuthSettings schema
type AuthSettings struct {
	AuthType interface{} `json:"authType"`
	CustomSettings interface{} `json:"customSettings,omitempty"`
	OAuth2Settings interface{} `json:"oAuth2Settings,omitempty"`
}
