// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OidcUserInfoEndpoint represents the OidcUserInfoEndpoint schema
// Endpoint for getting identity information about the user. For more information on the `/userinfo` endpoint, see [OpenID Connect](https://openid.net/specs/openid-connect-core-1_0.html#UserInfo).
type OidcUserInfoEndpoint struct {
	Binding interface{} `json:"binding,omitempty"`
	// URL of the resource server's `/userinfo` endpoint
	Url string `json:"url,omitempty"`
}
