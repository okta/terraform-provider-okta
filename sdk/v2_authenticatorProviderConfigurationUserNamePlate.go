// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AuthenticatorProviderConfigurationUserNamePlate struct {
	Template string `json:"template,omitempty"`
}

// Apple Push Notification Service
type APNS struct {
	ID               string `json:"id,omitempty"`
	AppBundleID      string `json:"appBundleId,omitempty"`
	DebugAppBundleID string `json:"debugAppBundleId,omitempty"`
}

// Firebase Cloud Messaging Service
type FCM struct {
	ID string `json:"id,omitempty"`
}
