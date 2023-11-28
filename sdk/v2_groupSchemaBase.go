// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type GroupSchemaBase struct {
	Id         string                           `json:"id,omitempty"`
	Properties map[string]*GroupSchemaAttribute `json:"properties,omitempty"`
	Required   []string                         `json:"required,omitempty"`
	Type       string                           `json:"type,omitempty"`
}
