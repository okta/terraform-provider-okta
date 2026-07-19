// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProtocolOidc represents the ProtocolOidc schema
// Protocol settings for authentication using the [OpenID Connect Protocol](http://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth)
type ProtocolOidc struct {
	// OpenID Connect and IdP-defined permission bundles to request delegated access from the user > **Note:** The [IdP type](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/Ide...
	Scopes []string `json:"scopes,omitempty"`
	Settings interface{} `json:"settings,omitempty"`
	// OpenID Connect Authorization Code flow
	Type string `json:"type,omitempty"`
	Algorithms interface{} `json:"algorithms,omitempty"`
	Credentials interface{} `json:"credentials,omitempty"`
	Endpoints interface{} `json:"endpoints,omitempty"`
	// URL of the IdP org
	OktaIdpOrgUrl string `json:"oktaIdpOrgUrl,omitempty"`
}
