// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// NetworkZoneAddress represents the NetworkZoneAddress schema
// Specifies the value of an IP address expressed using either `range` or `CIDR` form.
type NetworkZoneAddress struct {
	Type interface{} `json:"type,omitempty"`
	// Value in CIDR/range form, depending on the `type` specified
	Value string `json:"value,omitempty"`
}
