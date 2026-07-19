// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OSVersion represents the OSVersion schema
// Specifies the OS requirement for the policy.  There are two types of OS requirements:  * **Static**: A specific OS version requirement that doesn't change until you update the policy. A static OS r...
type OSVersion struct {
	// Contains the necessary properties for a dynamic version requirement
	DynamicVersionRequirement map[string]interface{} `json:"dynamicVersionRequirement,omitempty"`
	// The device version must be equal to or newer than the specified version string (maximum of three components for iOS and macOS, and maximum of four components for Android)
	Minimum string `json:"minimum,omitempty"`
}
