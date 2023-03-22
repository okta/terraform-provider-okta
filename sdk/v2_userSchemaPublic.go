package sdk

type UserSchemaPublic struct {
	Id         string                          `json:"id,omitempty"`
	Properties map[string]*UserSchemaAttribute `json:"properties,omitempty"`
	Required   []string                        `json:"required,omitempty"`
	Type       string                          `json:"type,omitempty"`
}
