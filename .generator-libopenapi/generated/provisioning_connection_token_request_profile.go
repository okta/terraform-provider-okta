// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProvisioningConnectionTokenRequestProfile represents the ProvisioningConnectionTokenRequestProfile schema
type ProvisioningConnectionTokenRequestProfile struct {
	AuthScheme interface{} `json:"authScheme"`
	// Token used to authenticate with the app
	Token string `json:"token,omitempty"`
}
