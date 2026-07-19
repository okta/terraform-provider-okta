// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OSVersionConstraint represents the OSVersionConstraint schema
type OSVersionConstraint struct {
	// The Windows device version must be equal to or newer than the specified version
	Minimum string `json:"minimum,omitempty"`
	// Contains the necessary properties for a dynamic Windows version requirement
	DynamicVersionRequirement map[string]interface{} `json:"dynamicVersionRequirement,omitempty"`
	// Indicates the Windows major version
	MajorVersionConstraint string `json:"majorVersionConstraint"`
}
