// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ServiceAccountDetailsAppAccountSub represents the ServiceAccountDetailsAppAccountSub schema
// Details for a SaaS app account, which will be managed as a service account
type ServiceAccountDetailsAppAccountSub struct {
	// The Okta app instance ID of the SaaS app
	OktaApplicationId string `json:"oktaApplicationId"`
	// The name of the SaaS app in the Okta Integration Network catalog
	AppGlobalName string `json:"appGlobalName,omitempty"`
	// The instance name of the SaaS app
	AppInstanceName string `json:"appInstanceName,omitempty"`
	Credentials interface{} `json:"credentials"`
}
