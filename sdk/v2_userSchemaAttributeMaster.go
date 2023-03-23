package sdk

type UserSchemaAttributeMaster struct {
	Priority []*UserSchemaAttributeMasterPriority `json:"priority,omitempty"`
	Type     string                               `json:"type,omitempty"`
}
