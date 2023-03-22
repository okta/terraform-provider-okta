package sdk

type InlineHookChannel struct {
	Config  *InlineHookChannelConfig `json:"config,omitempty"`
	Type    string                   `json:"type,omitempty"`
	Version string                   `json:"version,omitempty"`
}
