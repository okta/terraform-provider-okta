// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProvisioningConnectionOauthRequestProfile represents the ProvisioningConnectionOauthRequestProfile schema
type ProvisioningConnectionOauthRequestProfile struct {
	AuthScheme interface{} `json:"authScheme"`
	// Only used for the Okta Org2Org (`okta_org2org`) app. The unique client identifier for the OAuth 2.0 service app from the target org.
	ClientId string `json:"clientId,omitempty"`
	Settings interface{} `json:"settings,omitempty"`
	Signing interface{} `json:"signing,omitempty"`
}
