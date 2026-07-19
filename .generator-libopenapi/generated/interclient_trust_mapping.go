// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InterclientTrustMapping represents the InterclientTrustMapping schema
type InterclientTrustMapping struct {
	// ID of the org
	OrgId string `json:"orgId,omitempty"`
	// The app ID of the allowed app
	TrustedAppInstanceId string `json:"trustedAppInstanceId,omitempty"`
	// The app ID of the target app
	AppInstanceId string `json:"appInstanceId,omitempty"`
	// Timestamp when the interclient trust mapping was created
	Created string `json:"created,omitempty"`
	// The unique ID of the interclient trust mapping
	ID string `json:"id,omitempty"`
	// Timestamp when the interclient trust mapping was updated
	LastUpdated string `json:"lastUpdated,omitempty"`
	// ID of the user who created the interclient trust mapping
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
}
