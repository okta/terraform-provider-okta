package sdk

type EventHookChannel struct {
	Config  *EventHookChannelConfig `json:"config,omitempty"`
	Type    string                  `json:"type,omitempty"`
	Version string                  `json:"version,omitempty"`
}
