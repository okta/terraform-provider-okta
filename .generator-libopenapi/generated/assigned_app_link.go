// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AssignedAppLink represents the AssignedAppLink schema
type AssignedAppLink struct {
	AppAssignmentId string `json:"appAssignmentId,omitempty"`
	AppInstanceId string `json:"appInstanceId,omitempty"`
	AppName string `json:"appName,omitempty"`
	Hidden bool `json:"hidden,omitempty"`
	LinkUrl string `json:"linkUrl,omitempty"`
	SortOrder int `json:"sortOrder,omitempty"`
	CredentialsSetup bool `json:"credentialsSetup,omitempty"`
	ID string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
	LogoUrl string `json:"logoUrl,omitempty"`
}
