// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// DeviceContextProvider represents the DeviceContextProvider schema
type DeviceContextProvider struct {
	// Unique identifier for the device context provider
	ID string `json:"id,omitempty"`
	// Identifies the type of device context provider
	Key string `json:"key"`
	// Whether or not the device context provider is used to identify the user. `IGNORE` prevents the device context provider from being used to authenticate the user. Identification of the device and dev...
	UserIdentification string `json:"userIdentification,omitempty"`
}
