// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type Scope struct {
	AllowedOktaApps []*IframeEmbedScopeAllowedApps `json:"allowedOktaApps,omitempty"`
	StringValue     string                         `json:"stringValue,omitempty"`
	Type            string                         `json:"type,omitempty"`
}
