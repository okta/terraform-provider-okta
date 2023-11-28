// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type EventHookChannel struct {
	Config  *EventHookChannelConfig `json:"config,omitempty"`
	Type    string                  `json:"type,omitempty"`
	Version string                  `json:"version,omitempty"`
}
