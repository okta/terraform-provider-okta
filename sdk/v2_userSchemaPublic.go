// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserSchemaPublic struct {
	Id         string                          `json:"id,omitempty"`
	Properties map[string]*UserSchemaAttribute `json:"properties,omitempty"`
	Required   []string                        `json:"required,omitempty"`
	Type       string                          `json:"type,omitempty"`
}
