// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PerClientRateLimitSettings represents the PerClientRateLimitSettings schema
type PerClientRateLimitSettings struct {
	// The default PerClientRateLimitMode that applies to any use case in the absence of a more specific override
	DefaultMode interface{} `json:"defaultMode"`
	// A map of Per-Client Rate Limit Use Case to the applicable PerClientRateLimitMode. Overrides the `defaultMode` property for the specified use cases.
	UseCaseModeOverrides map[string]interface{} `json:"useCaseModeOverrides,omitempty"`
}
