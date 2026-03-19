// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserSchemaAttributeItems struct {
	Enum  []interface{}              `json:"enum,omitempty"`
	OneOf []*UserSchemaAttributeEnum `json:"oneOf,omitempty"`
	Type  string                     `json:"type,omitempty"`
}
