// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CreateGroupPushMappingRequest represents the CreateGroupPushMappingRequest schema
type CreateGroupPushMappingRequest struct {
	// The ID of the existing target group for the group push mapping. This is used to link to an existing group. Required if `targetGroupName` is not provided.
	TargetGroupId string `json:"targetGroupId,omitempty"`
	// The name of the target group for the group push mapping. This is used when creating a new downstream group. If the group already exists, it links to the existing group. Required if `targetGroupId` ...
	TargetGroupName string `json:"targetGroupName,omitempty"`
	AppConfig interface{} `json:"appConfig,omitempty"`
	// The ID of the source group for the group push mapping
	SourceGroupId string `json:"sourceGroupId"`
	Status string `json:"status,omitempty"`
}
