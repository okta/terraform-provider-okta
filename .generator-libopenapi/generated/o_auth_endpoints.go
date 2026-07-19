// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuthEndpoints represents the OAuthEndpoints schema
// The `OAUTH2` and `OIDC` protocols support the `authorization` and `token` endpoints. Also, the `OIDC` protocol supports the `userInfo` and `jwks` endpoints.  The IdP Authorization Server (AS) endpo...
type OAuthEndpoints struct {
	UserInfo interface{} `json:"userInfo,omitempty"`
	Authorization interface{} `json:"authorization,omitempty"`
	Jwks interface{} `json:"jwks,omitempty"`
	Slo interface{} `json:"slo,omitempty"`
	Token interface{} `json:"token,omitempty"`
}
