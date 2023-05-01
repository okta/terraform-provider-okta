package sdk

type UserSchemaAttributeItems struct {
	Enum  []interface{}              `json:"enum,omitempty"`
	OneOf []*UserSchemaAttributeEnum `json:"oneOf,omitempty"`
	Type  string                     `json:"type,omitempty"`
}
