// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppAccountContainerDetails represents the AppAccountContainerDetails schema
// Container details for resource type APP_ACCOUNT
type AppAccountContainerDetails struct {
	// Human-readable name of the container that owns the privileged resource
	DisplayName string `json:"displayName,omitempty"`
	// The application global ID
	GlobalAppId string `json:"globalAppId,omitempty"`
	// Indicates if the application supports password push
	PasswordPushSupported bool `json:"passwordPushSupported,omitempty"`
	// Indicates if provisioning is enabled for this application
	ProvisioningEnabled bool `json:"provisioningEnabled,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// The application name
	AppName string `json:"appName,omitempty"`
	// The app ID associated with the privileged resource
	ContainerId string `json:"containerId"`
}
