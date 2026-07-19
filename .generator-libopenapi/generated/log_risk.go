// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogRisk represents the LogRisk schema
// Risk associated with the event
type LogRisk struct {
	// The name of the detection mechanism that identified the risk
	DetectionName string `json:"detectionName,omitempty"`
	// The entity that issued the associated risk
	Issuer string `json:"issuer,omitempty"`
	// The risk level associated with the request
	Level string `json:"level,omitempty"`
	// The previous risk level (if any) associated with the user
	PreviousLevel string `json:"previousLevel,omitempty"`
	// Reasons for the associated risk
	Reasons []string `json:"reasons,omitempty"`
}
