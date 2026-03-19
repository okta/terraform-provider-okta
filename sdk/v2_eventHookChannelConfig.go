// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type EventHookChannelConfig struct {
	AuthScheme *EventHookChannelConfigAuthScheme `json:"authScheme,omitempty"`
	Headers    []*EventHookChannelConfigHeader   `json:"headers,omitempty"`
	Uri        string                            `json:"uri,omitempty"`
}
