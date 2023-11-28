// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ApplicationAccessibility struct {
	ErrorRedirectUrl string `json:"errorRedirectUrl,omitempty"`
	LoginRedirectUrl string `json:"loginRedirectUrl,omitempty"`
	SelfService      *bool  `json:"selfService,omitempty"`
}
