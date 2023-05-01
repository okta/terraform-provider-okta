package sdk

type EventHookChannelConfig struct {
	AuthScheme *EventHookChannelConfigAuthScheme `json:"authScheme,omitempty"`
	Headers    []*EventHookChannelConfigHeader   `json:"headers,omitempty"`
	Uri        string                            `json:"uri,omitempty"`
}
