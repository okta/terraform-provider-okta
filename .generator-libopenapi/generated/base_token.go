// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BaseToken represents the BaseToken schema
type BaseToken struct {
	// Claims included in the token. Consists of name-value pairs for each included claim. For descriptions of the claims that you can include, see the Okta [OpenID Connect and OAuth 2.0 API reference](/o...
	Claims map[string]interface{} `json:"claims,omitempty"`
	// The token
	Token map[string]interface{} `json:"token,omitempty"`
}
