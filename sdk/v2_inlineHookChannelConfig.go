// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type InlineHookChannelConfig struct {
	AuthScheme   *InlineHookChannelConfigAuthScheme `json:"authScheme,omitempty"`
	Headers      []*InlineHookChannelConfigHeaders  `json:"headers,omitempty"`
	Method       string                             `json:"method,omitempty"`
	Uri          string                             `json:"uri,omitempty"`
	AuthType     string                             `json:"authType,omitempty"`
	ClientId     string                             `json:"clientId,omitempty"`
	ClientSecret string                             `json:"clientSecret,omitempty"`
	TokenUrl     string                             `json:"tokenUrl,omitempty"`
	Scope        string                             `json:"scope,omitempty"`
}
