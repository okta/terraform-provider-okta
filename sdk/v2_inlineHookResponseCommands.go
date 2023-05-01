package sdk

type InlineHookResponseCommands struct {
	Type  string                            `json:"type,omitempty"`
	Value []*InlineHookResponseCommandValue `json:"value,omitempty"`
}
