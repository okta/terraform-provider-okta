// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type InlineHookChannelConfig struct {
	AuthScheme *InlineHookChannelConfigAuthScheme `json:"authScheme,omitempty"`
	Headers    []*InlineHookChannelConfigHeaders  `json:"headers,omitempty"`
	Method     string                             `json:"method,omitempty"`
	Uri        string                             `json:"uri,omitempty"`
}
