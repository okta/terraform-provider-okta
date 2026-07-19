// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ProvisioningConnectionProfileOauth represents the ProvisioningConnectionProfileOauth schema
// The app provisioning connection profile used to configure the method of authentication and the credentials. Currently, token-based and OAuth 2.0-based authentication are supported.
type ProvisioningConnectionProfileOauth struct {
	AuthScheme interface{} `json:"authScheme"`
	ClientId string `json:"clientId,omitempty"`
}
