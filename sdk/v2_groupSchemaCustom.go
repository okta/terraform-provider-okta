package sdk

type GroupSchemaCustom struct {
	Id         string                           `json:"id,omitempty"`
	Properties map[string]*GroupSchemaAttribute `json:"properties,omitempty"`
	Required   []string                         `json:"required,omitempty"`
	Type       string                           `json:"type,omitempty"`
}
