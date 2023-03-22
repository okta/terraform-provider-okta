package sdk

type InlineHookChannelConfig struct {
	AuthScheme *InlineHookChannelConfigAuthScheme `json:"authScheme,omitempty"`
	Headers    []*InlineHookChannelConfigHeaders  `json:"headers,omitempty"`
	Method     string                             `json:"method,omitempty"`
	Uri        string                             `json:"uri,omitempty"`
}
