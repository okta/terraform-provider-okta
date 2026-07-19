// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// WellKnownAppAuthenticatorConfiguration represents the WellKnownAppAuthenticatorConfiguration schema
type WellKnownAppAuthenticatorConfiguration struct {
	// The authenticator enrollment endpoint
	AppAuthenticatorEnrollEndpoint string `json:"appAuthenticatorEnrollEndpoint,omitempty"`
	// The unique identifier of the app authenticator
	AuthenticatorId string `json:"authenticatorId,omitempty"`
	// Timestamp when the authenticator was created
	CreatedDate *time.Time `json:"createdDate,omitempty"`
	Key interface{} `json:"key,omitempty"`
	// The authenticator display name
	Name string `json:"name,omitempty"`
	// The `id` of the Okta Org
	OrgId string `json:"orgId,omitempty"`
	SupportedMethods []interface{} `json:"supportedMethods,omitempty"`
	// Timestamp when the authenticator was last modified
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	// The type of authenticator
	Type string `json:"type,omitempty"`
}
